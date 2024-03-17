#include <stdio.h>
#include <stdlib.h>

int main() {
	while (1) {
        char *cmd = NULL;
        size_t size = 0;
        
        printf("gsh~> ");

		if (getline(&cmd, &size, stdin) == -1) {
			puts("Error reading input");
			exit(EXIT_FAILURE);
		}

		printf("got command %s", cmd);

		free(cmd);
	}

	return 0;
}
