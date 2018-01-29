// TODO(olegs): Technically, we could completely avoid having this file
// and write everything in bindings.pyx, but the code generated from it
// is very difficult to debug, so I'd rather move more code from *.pyx
// into *.c even though it requires writing in one extra language.
#include "libgokubectl.h"
#include <stdio.h>
#include <stddef.h>
#include <stdlib.h>
#include <string.h>

struct ResourceGet_return kubectl_get(
    const char* opts,
    size_t opts_len,
    const char** args,
    size_t nargs) {
    GoString* items = (GoString*)malloc(nargs * sizeof(GoString));
    size_t i = 0;

    while (i < nargs) {
        items[i].p = args[i];
        items[i].n = (GoInt)strlen(args[i]);
        i++;
    }
    GoSlice pods_slice = {
        .data = items,
        .len = (GoInt)nargs,
        .cap = (GoInt)nargs,
    };
    GoString json_opts = { .p = opts, .n = opts_len };
    struct ResourceGet_return result = ResourceGet(json_opts, pods_slice);
    free(items);
    return result;
}

struct Create_return kubectl_create(const char* opts, size_t opts_len) {
    GoString items = { .p = opts, .n = opts_len };
    return Create(items);
}


void free_cstring_array(char** strings, size_t n) {
    size_t i = 0;
    while (i < n) {
        free(strings[i]);
        i++;
    }
    free(strings);
}

void free_resource_get(struct ResourceGet_return rgr) {
    free(rgr.r0);
    free(rgr.r1);
}

void free_create(struct Create_return cr) {
    free(cr.r0);
    free(cr.r1);
}
