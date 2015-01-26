#include "miharihook.h"
#include <stdio.h>
#include <dlfcn.h>
#include <stdlib.h>
#include <unistd.h>
#include <string.h>
#include <sys/param.h>
#include <stdarg.h>
#include <sys/stat.h>
#include <errno.h>

static int enable_debug;
static void *libc_pointer;

static FILE *(*libc_fopen)(const char *restrict filename, const char *restrict mode);
static int (*libc_open)(const char *path, int oflag, ...);
static int (*libc_openat)(int fd, const char *path, int oflag, ...);

static int logfd;

static void close_log(void)
{
    close(logfd);
}

static void open_log(void)
{
#ifdef __APPLE__
    libc_pointer = dlopen("/usr/lib/libSystem.dylib", RTLD_LAZY);
#elif defined(__linux__)
    libc_pointer = dlopen("/lib/libc.so.6", RTLD_LAZY);
#else
#error "Unknown platform"
#endif
    if (libc_pointer == NULL) {
        fprintf(stderr, "Cannot open C library");
        abort();
    }

    if (getenv("MIHARI_DEBUG")) {
        enable_debug = 1;
    }

    libc_open = dlsym(libc_pointer, "open");
    libc_openat = dlsym(libc_pointer, "openat");
    libc_fopen = dlsym(libc_pointer, "fopen");
    
    logfd = libc_open(getenv("MIHARI_LOG_FILE"), O_WRONLY|O_APPEND);
    if (logfd < 0)
        perror("CANNOT OPEN MIHARI LOG FILE");
    atexit(close_log);
}

static void log_open(const char *path, char flag)
{
    if (logfd > 0) {
        char buf[MAXPATHLEN*2+3];
        bzero(buf, sizeof(buf));
        char wd[MAXPATHLEN];
        getcwd(wd, sizeof(wd));
        
        size_t length = snprintf(buf, sizeof(buf)-2, "O%c%s:%s", flag, wd, path);
        if (enable_debug)
            fprintf(stderr, "%s\n", buf);
        write(logfd, buf, length+1);
    }
}

static int is_file_exists(const char *path)
{
    struct stat buf;
    int ret = stat(path, &buf);
    (void)buf.st_uid;
    if (ret == 0) return 1;
    switch (errno) {
    case EACCES:
        return 1;
    default:
        return 0;
    }
}

static char open_flag_parse(const char *path, int oflag)
{
    char flag;
    switch (oflag & O_ACCMODE) {
    case O_RDONLY:
        flag = 'R';
        break;
    case O_WRONLY:
        flag = 'W';
        break;
    case O_RDWR:
        if (oflag & O_TRUNC) {
            flag = 'W';
        } else if (oflag & O_CREAT) {
            struct stat buf;
            if (!is_file_exists(path))
                flag = 'W';
            else
                flag = 'B';
        } else {
            flag = 'B';
        }
        break;
    default:
        flag = 'U';
        break;
    }
    return flag;
}

int
open(const char *path, int oflag, ...)
{
    if (libc_pointer == NULL)
        open_log();
    log_open(path, open_flag_parse(path, oflag));

    va_list a;
	va_start(a, oflag);
    mode_t mode;
	mode = va_arg(a, int);
	va_end(a);

    return libc_open(path, oflag, mode);
}

int
openat(int fd, const char *path, int oflag, ...)
{
    if (libc_pointer == NULL)
        open_log();
    log_open(path, open_flag_parse(path, oflag));

    va_list a;
	va_start(a, oflag);
    mode_t mode;
	mode = va_arg(a, int);
	va_end(a);

    return libc_openat(fd, path, oflag, mode);
}

static char fopen_flag(const char *restrict filename, const char *restrict mode)
{
    int read_mode = 0;
    int write_mode = 0;
    char flag = 'U';

    if (strcmp(mode, "r")) read_mode = 1;
    else if (strcmp(mode, "r+")) {read_mode = 1;}
    else if (strcmp(mode, "w")) write_mode = 1;
    else if (strcmp(mode, "w+")) write_mode = 1;
    else if (strcmp(mode, "a")) {read_mode = 1; write_mode = 1;}
    else if (strcmp(mode, "a+")) {read_mode = 1; write_mode = 1;}
    
    if (write_mode && read_mode) {
        if (is_file_exists(filename)) {
            flag = 'B';
        } else {
            flag = 'W';
        }
    } else if (write_mode) {
        flag = 'W';
    } else if (read_mode) {
        flag = 'R';
    }
    return flag;
}

FILE *
fopen(const char *restrict filename, const char *restrict mode)
{
    if (libc_pointer == NULL)
        open_log();
    log_open(filename, fopen_flag(filename, mode));

    return libc_fopen(filename, mode);
}
