# -*- coding: utf-8 -*-

from libc.stdio cimport printf
from libc.stdlib cimport malloc, free
from libc.string cimport strcmp, memcpy


cdef extern from "lib/bindings.h":
    struct GoString:
        const char* p
        int n

    struct ResourceGet_return:
        GoString r0
        GoString r1

    struct Create_return:
        GoString r0
        GoString r1

    ResourceGet_return kubectl_get(
        const char* opts,
        size_t opts_len,
        const char** args,
        size_t nargs
    )

    Create_return kubectl_create(const char* opts, size_t opts_len)

    void free_cstring_array(char**, size_t)


cdef char** to_cstring_array(object strings):
    cdef char** result = <char**>malloc(len(strings) * sizeof(char*))
    cdef char* itemp
    cdef size_t string_size
    cdef bytes item
    cdef char* citem
    cdef size_t nelements = len(strings)
    cdef size_t i = 0

    while i < nelements:
        item = str(strings[i]).encode("utf-8")
        itemp = <char*>item
        string_size = <size_t>(len(item) + 1) * sizeof(char)
        citem = <char*>malloc(string_size)
        memcpy(citem, itemp, string_size)
        citem[string_size - 1] = <char>0
        result[i] = citem
        i += 1
    return result


cpdef object pykubectl_get_impl(object items, bytes options):
    cdef size_t nargs = len(items)
    cdef size_t optlen = len(options)
    cdef bytes message
    cdef char** args = to_cstring_array(items)
    cdef char* opts = options

    cdef ResourceGet_return result = kubectl_get(
        opts,
        optlen,
        <const char**>args,
        nargs
    )
    free_cstring_array(args, nargs)

    if result.r0.n == 0:
        message = result.r1.p[:result.r1.n]
        raise Exception("kubectl failed: '{}'".format(message.decode("utf-8")))
    message = result.r0.p[:result.r0.n]
    return message.decode("utf-8")


cpdef object pykubectl_create_impl(bytes options):
    cdef const char* opts = options
    cdef size_t opt_len = len(options)
    cdef Create_return result = kubectl_create(opts, opt_len)
    cdef bytes message

    if result.r0.n == 0:
        message = result.r1.p[:result.r1.n]
        raise Exception("kubectl failed: '{}'".format(message.decode("utf-8")))
    message = result.r0.p[:result.r0.n]
    return message.decode("utf-8")
