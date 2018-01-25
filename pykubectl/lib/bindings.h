#ifndef BINDINGS_H
#define BINDINGS_H

#include "libgokubectl.h"
#include <stddef.h>

struct ResourceGet_return kubectl_get(const char**, size_t);

#endif // BINDINGS_H
