#!/usr/bin/env python3.6
# -*- coding: utf-8 -*-

import os
import re
import shutil
import subprocess

from setuptools import setup, Extension
from Cython.Distutils import build_ext
from setuptools.command.build_py import build_py


LIB_DIR = '/usr/local/lib'


class BuildGo(build_py):
    '''
    Build Go bindings before building Python extension
    '''

    def run(self):
        version = None
        refusal = 'Will not compile Go bindings'
        try:
            version = subprocess.check_output(['go', 'version'])
        except subprocess.CalledProcessError:
            print('Didn\'t find Go compiler.')
            print(refusal)
            return
        match = re.search(br'go(\d+)\.(\d+)\.(\d+)', version)
        if match:
            version = tuple(map(int, match.group(1, 2, 3)))
        else:
            print('Unrecognized version of Go compiler: {}'.format(version))
            print(refusal)
            return
        req_version = 1, 9, 1
        if version < req_version:
            print('You need Go compiler version higher than: {}'.format(
                req_version,
            ))
            print(refusal)
            return
        try:
            result = subprocess.check_output([
                'go',
                'build',
                '-v',
                '-buildmode=c-shared',
                '-o',
                'pykubectl/lib/libgokubectl.so',
                'github.com/wvxvw-traiana/pykubectl/main',
            ])
            print(result)
        except subprocess.CalledProcessError as e:
            print(e.stderr)
            raise
        build_py.run(self)


def is_extension_file(name):
    return name.startswith("pykubectl.wrapped") and not (
        name.endswith(".pyx") or name.endswith(".pxd")
    )


for root, dirs, files in os.walk(".", topdown=False):
    for name in files:
        if is_extension_file(name):
            os.remove(os.path.join(root, name))
    for name in dirs:
        if name == "build":
            shutil.rmtree(name)

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
