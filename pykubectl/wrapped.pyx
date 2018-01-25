# -*- coding: utf-8 -*-

import json

from libc.stdio cimport FILE, fopen, fwrite, fclose, stdout
from libc.stdlib cimport malloc, free
from libc.string cimport strcmp


cdef extern from "lib/bindings.h":
    struct GoString:
        const char* p
        int n

    struct ResourceGet_return:
        GoString r0
        GoString r1

    ResourceGet_return kubectl_get(const char** args, size_t nargs)


cdef const char** to_cstring_array(strings):
    cdef char** result = <char**>malloc(len(strings) * sizeof(char*))
    cdef bytes item
    for i in xrange(len(strings)):
        s = str(strings[i]).encode("utf-8")
        item = s
        result[i] = item
    return result


cpdef object pykubectl_get(object options, object items):
    cdef size_t nargs = len(items);
    cdef const char** args = to_cstring_array(items);
    cdef ResourceGet_return result = kubectl_get(args, nargs);
    cdef bytes message

    if result.r0.n == 0:
        message = result.r1.p
        raise Exception("kubectl failed: '{}'".format(message.decode("utf-8")))
    message = result.r0.p
    return json.loads(message.decode("utf-8"))
