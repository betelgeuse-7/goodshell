#include <stdlib.h>
#include <stdio.h>
#include <string.h>
#include "runcmd.h"

#define PATH_DELIM ":"

void run_cmd(char *cmd) {
	char *exe_loc = get_exe_loc(cmd);
}

char *get_exe_loc(char *cmd) {
	char *path = getenv("PATH");
    char *token;

    token = strtok(path, PATH_DELIM);
    while (token != NULL) {
        printf("%s\n", token);
        token = strtok(NULL, PATH_DELIM);
    }
}