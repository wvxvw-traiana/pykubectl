#!/usr/bin/env python3.6
# -*- coding: utf-8 -*-

import os
import re
import shutil
import subprocess
import pathlib

from Cython.Distutils import build_ext
from setuptools import setup, Extension
from setuptools.command.build_py import build_py
from distutils.dep_util import newer_group


LIB_DIR = '/usr/local/lib'


class BuildGo(build_py):
    '''
    Build Go bindings before building Python extension
    '''

    gokubectl_package = 'github.com/wvxvw-traiana/pykubectl/main'

    def run(self):
        if self.stale_go():
            print('Fresh Go build')
            self.build_go()
        else:
            print('Using previously compiled Go library')
        if self.stale_win_go():
            print('Fresh Win Go build')
            self.build_win_go()
        build_py.run(self)

    def stale_go(self):
        installed_lib = os.path.join(LIB_DIR, 'libgokubectl.so')
        go_package = pathlib.Path('main')
        go_sources = list(str(f.resolve()) for f in go_package.rglob('*.go'))
        return newer_group(go_sources, installed_lib)

    def stale_win_go(self):
        existing_archive = 'pykubectl/lib/libgokubectl.dll'
        go_package = pathlib.Path('main')
        go_sources = list(str(f.resolve()) for f in go_package.rglob('*.go'))
        return newer_group(go_sources, existing_archive)

    def is_go_available(self):
        version = None
        refusal = 'Will not compile Go bindings'
        try:
            version = subprocess.check_output(['go', 'version'])
        except subprocess.CalledProcessError:
            print('Didn\'t find Go compiler.')
            print(refusal)
            return False
        match = re.search(br'go(\d+)\.(\d+)\.(\d+)', version)
        if match:
            version = tuple(map(int, match.group(1, 2, 3)))
        else:
            print('Unrecognized version of Go compiler: {}'.format(version))
            print(refusal)
            return False
        req_version = 1, 9, 1
        if version < req_version:
            print('You need Go compiler version higher than: {}'.format(
                req_version,
            ))
            print(refusal)
            return False
        return True

    def is_gcc_available(self):
        try:
            subprocess.check_output([
                'x86_64-w64-mingw32-gcc-win32', '--version'
            ])
        except subprocess.CalledProcessError:
            print('Didn\'t find GCC')
            print('Will not compile Windows DLL')
            return False
        return True

    def build_win_go(self):
        if self.is_go_available() and self.is_gcc_available():
            try:
                args = [
                    'go',
                    'build',
                    '-v',
                    '-buildmode=c-archive',
                    '-o',
                    'pykubectl/lib/libgokubectl.a',
                    self.gokubectl_package,
                ]
                env = dict(os.environ)
                env.update({
                    'GOOS': 'windows',
                    'GOARCH': 'amd64',
                    'CGO_ENABLED': '1',
                    'CC': 'x86_64-w64-mingw32-gcc-win32',
                })
                proc = subprocess.Popen(args, env=env)
                stdout, stderr = proc.communicate()
                print(stdout)
                if proc.returncode != 0:
                    raise Exception('Go compilation failed:\n{}'.format(
                        stderr
                    ))
            except subprocess.CalledProcessError as e:
                print(e.stderr)
                raise
            try:
                result = subprocess.check_output([
                    'x86_64-w64-mingw32-gcc-win32',
                    '-shared',
                    '-pthread',
                    '-o',
                    'pykubectl/lib/gokubectl.dll',
                    'pykubectl/lib/win_gokubectl.c',
                    'pykubectl/lib/libgokubectl.a',
                    '-lwinmm',
                    '-lntdll',
                    '-lws2_32',
                ])
                print(result)
            except subprocess.CalledProcessError as e:
                print(e.stderr)
                raise

    def build_go(self):
        if self.is_go_available():
            try:
                result = subprocess.check_output([
                    'go',
                    'build',
                    '-v',
                    '-buildmode=c-shared',
                    '-o',
                    'pykubectl/lib/libgokubectl.so',
                    self.gokubectl_package,
                ])
                print(result)
            except subprocess.CalledProcessError as e:
                print(e.stderr)
                raise


# def is_extension_file(name):
#     return name.startswith("pykubectl.wrapped") and not (
#         name.endswith(".pyx") or name.endswith(".pxd")
#     )


# TODO(olegs): This should be smarter to only delete stale binaries
# for root, dirs, files in os.walk(".", topdown=False):
#     for name in files:
#         if is_extension_file(name):
#             os.remove(os.path.join(root, name))
#     for name in dirs:
#         if name == "build":
#             shutil.rmtree(name)

setup(
    install_requires=['cython'],
    packages=['pykubectl'],
    zip_safe=False,
    name='pykubectl',
    version='0.0.1',
    description='Python binding for kubectl',
    author='Oleg Sivokon',
    author_email='olegs@traiana.com',
    url='TBD',
    license='PROPRIETARY',
    cmdclass={
        'build_ext': build_ext,
        'build_py': BuildGo,
    },
    scripts=['bin/pykubectl'],
    # See: http://blog.codekills.net/2011/07/15/\
    # lies,-more-lies-and-python-packaging-documentation-on--package_data-/
    # package_data={'pykubectl': ['lib/libgokubectl.so']},
    data_files=[(LIB_DIR, ['pykubectl/lib/libgokubectl.so'])],
    ext_modules=[
        Extension(
            'pykubectl.wrapped',
            sources=[
                'pykubectl/lib/bindings.c',
                'pykubectl/wrapped.pyx',
            ],
            include_dirs=['pykubectl/lib'],
            library_dirs=['pykubectl/lib'],
            libraries=['gokubectl'],
            extra_link_args=['-Wl,-rpath,{}'.format(LIB_DIR)],
            # runtime_library_dirs=[
            #     '/usr/local/lib/python3.6/site-packages/pykubectl-0.0.1-py3.6-linux-x86_64.egg/pykubectl/lib',
            #     './lib', './pykubectl/lib', '/usr/local/lib',
            # ],
        )
    ]
)
