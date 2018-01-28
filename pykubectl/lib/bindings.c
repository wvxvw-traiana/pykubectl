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

    printf("C items: %s\n", args[0]);
    printf("C opts: %s\n", opts);
    
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


void free_cstring_array(char** strings, size_t n) {
    size_t i = 0;
    while (i < n) {
        free(strings[i]);
        i++;
    }
    free(strings);
}
