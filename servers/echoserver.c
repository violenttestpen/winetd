#include <stdio.h>

int main() {
    int n = 1;
    char buf[1];
    while (n > 0) {
        n = fread(buf, 1, sizeof(buf), stdin);
        fwrite(buf, 1, sizeof(buf), stdout);
    }
    return 0;
}