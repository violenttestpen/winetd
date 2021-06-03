#include <stdio.h>
#include <stdlib.h>

void ignore_me_init_buffering() {
    setvbuf(stdin, NULL, _IONBF, 0);
    setvbuf(stdout, NULL, _IONBF, 0);
    setvbuf(stderr, NULL, _IONBF, 0);
}

void win() {
    system("whoami");
    puts("Sorry, were you expecting a shell? Here, let me give you cmd.exe");
    exit(0);
}

void welcome() {
    char buf[0x40];
    printf("Enter your name:\n");
    gets(buf);
    printf("Welcome, ");
    printf(buf);
    putchar('\n');
}

void vuln() {
    char buf[0x40];
    printf("Enter your input:\n");
    gets(buf);
}

void main() {
    ignore_me_init_buffering();
    welcome();
    vuln();
}


