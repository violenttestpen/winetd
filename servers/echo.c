#include <stdio.h>

void ignore_me_init_buffering() {
    setvbuf(stdin, NULL, _IONBF, 0);
    setvbuf(stdout, NULL, _IONBF, 0);
    setvbuf(stderr, NULL, _IONBF, 0);
}

int main() {
    ignore_me_init_buffering();
    char buf[128];
    printf("buf: %p\n", buf);
    printf("Enter input: ");
    fgets(buf, 0x200, stdin);
    printf(buf);
    return 0;
}