#ifndef BINDINGS_H
#define BINDINGS_H

#include "libgokubectl.h"
#include <stddef.h>

struct ResourceGet_return kubectl_get(const char*, size_t, const char**, size_t);

struct Create_return kubectl_create(const char*, size_t);

void free_cstring_array(char**, size_t);

void free_resource_get(struct ResourceGet_return);

void free_create(struct Create_return);

#endif // BINDINGS_H
