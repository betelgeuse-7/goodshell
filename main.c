#include <stdio.h>
#include <stdlib.h>

#include "ParseTree.h";

int process_command(char *cmd);

int main(void) {
    while (1)
    {
        puts("gsh~> ");
        char *line = "";
        size_t len = 0;
        size_t read;

        if ((read = getline(&line, &len, stdin)) == -1) {
            puts("Error reading the command");
            return 1;
        }

        if ((process_command(line)) == -1) {
            puts("Error handling the command");
            return 1;
        }
    }
}

int process_command(char *cmd) {
    ParseTree *tree = parse_command(cmd);
    // evaluate the tree
}
