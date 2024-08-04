#include <stdio.h>
#include <stdlib.h>
#include <stdbool.h>
#include <string.h>

#include "Buffer.h"

bool process_line(char *cmd);
bool is_letter(char);

int main(void) {
    while (true)
    {
        printf("gsh~> ");
        char *line = "";
        size_t len = 0;
        size_t read;

        if ((read = getline(&line, &len, stdin)) == -1) {
            puts("Error reading the command");
            return 1;
        }

        if (!process_line(line)) {
            puts("Error handling the command");
            return 1;
        }
    }
}

bool process_line(char *line) {
    if (strlen(line) == 0) {
        return true;
    }

    if (is_letter(line[0])) {
        char *cmd = NULL;
        int i = 0;
        Buffer *buf = new_buffer();

        while (true)
        {
            if (!is_letter(line[i])) break;
            Buffer_append(buf, line[i]);
            i++;
        }

        Buffer_print(buf);
    }

    return true;
}

bool is_letter(char ch) {
    return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z');
}
