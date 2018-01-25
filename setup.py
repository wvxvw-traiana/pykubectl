#!/usr/bin/env python3.6
# -*- coding: utf-8 -*-

import os
import shutil

from setuptools import setup, Extension
from Cython.Distutils import build_ext


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
    cmdclass={"build_ext": build_ext},
    scripts=['bin/pykubectl'],
    # See: http://blog.codekills.net/2011/07/15/\
    # lies,-more-lies-and-python-packaging-documentation-on--package_data-/
    package_data={'pykubectl': ['lib/libgokubectl.so']},
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
            extra_link_args=['-Wl,-rpath,/usr/local/lib'],
            # -Wl,+s
            # runtime_library_dirs=[
            #     '/usr/local/lib/python3.6/site-packages/pykubectl-0.0.1-py3.6-linux-x86_64.egg/pykubectl/lib',
            #     './lib', './pykubectl/lib', '/usr/local/lib',
            # ],
        )
    ]
)
