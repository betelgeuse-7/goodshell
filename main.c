#include <stdio.h>
#include <stdlib.h>
#include <stdbool.h>

#include "runcmd.h"

int main() {
	while (true) {
        char *cmd = NULL;
        size_t size = 0;
        
        printf("gsh~> ");

		if (getline(&cmd, &size, stdin) == -1) {
			puts("Error reading input");
			exit(EXIT_FAILURE);
		}

		run_cmd(cmd);
		free(cmd);
	}

	return 0;
}
