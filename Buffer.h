#include <stdlib.h>
#include <stdio.h>

#define DEFAULT_BUFFER_SIZE 25

/**
 * A simple ad-hoc string buffer
 */
typedef struct {
    char *data;
    size_t len;
    size_t cap;
} Buffer;

Buffer *new_buffer(void) {
    Buffer *b = (Buffer *) malloc(sizeof(Buffer));
    b->data = (char *) malloc(DEFAULT_BUFFER_SIZE);
    if (b->data == NULL) {
        fprintf(stderr, "Memory allocation for new_buffer failed");
        exit(EXIT_FAILURE);
    }
    b->len = 0;
    b->cap = DEFAULT_BUFFER_SIZE;
    return b;
}

void Buffer_append(Buffer *b, char c) {
    if (b->len >= b->cap) {
        size_t new_cap = b->cap * 2;
        b->data = (char *)realloc(b->data, new_cap);
        if (b->data == NULL) {
            fprintf(stderr, "Memory reallocation for Buffer_append failed");
            exit(EXIT_FAILURE);
        }

        b->cap = new_cap;
    }

    b->data[b->len++] = c;
    b->data[b->len] = '\0';
}

char Buffer_get(Buffer *b, size_t idx) {
    if (idx >= b->len) {
        fprintf(stderr, "Buffer_get: index out of bounds: len=%ld, idx=%ld\n", b->len, idx);
        exit(EXIT_FAILURE);
    }

    return b->data[idx];
}

void Buffer_free(Buffer *b) {
    free(b->data);
    free(b);
}

char *Buffer_data(Buffer *b) {
    return b->data;
}

size_t Buffer_len(Buffer *b) {
    return b->len-1; // exclude the NUL terminator
}

size_t Buffer_cap(Buffer *b) {
    return b->cap;
}

void Buffer_print(Buffer *b) {
    if (b->len == 0) {
        printf("Buffer { data: , len: 1, cap: %d }\n", DEFAULT_BUFFER_SIZE);
        return;
    }

    printf("Buffer{\n\tdata: %s\n\tlen: %ld\n\tcap: %ld\n}\n", b->data, b->len, b->cap);
}