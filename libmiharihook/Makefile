# Hey Emacs, this is a -*- makefile -*-
# Copyright (C) 2013 Y.Okamura
# https://gist.github.com/2343121
# MIT License

# Target Executable file
TARGET = libmiharihook

# C source code
SRC = miharihook.c

# CPP source code
CXXSRC = 

# Common Flags
CPPFLAGS += -Wall -Wextra -pipe
COMMON_FLAGS += 

# C compiler
CC = gcc
CFLAGS +=

# C++ compiler
CXX = g++
CXXFLAGS +=

# AR
AR = ar
ARFLAGS = urcv

# Linker
LD = gcc
LDFLAGS += 
LIBS += -lstdc++ -lm

# Windows Resource Compiler
WINDRES = windres

# Unix like OS
UNAME = $(strip $(shell uname))
CPPFLAGS += 
LDFLAGS  +=
LIBS += 

ifeq ($(UNAME), Darwin) # Mac OS X
  LDFLAGS += -dynamiclib -fPIC
  LIB = $(TARGET).so
endif
ifeq ($(UNAME), Linux) # Linux
  LDFLAGS +=
  LIB = $(TARGET).so
endif



# For Emacs
ifdef EMACS
ifeq ($(CC),clang)
  CPPFLAGS += -fno-color-diagnostics
endif
endif

# Debug Options
ifdef RELEASE
  CPPFLAGS += -Os -mtune=native
else
  CPPFLAGS += -g -O0
endif

# copy commong flags
CPPFLAGS += $(COMMON_FLAGS)
LDFLAGS += $(COMMON_FLAGS)

# generate dependence file
OBJDIR = ./objs
CPPFLAGS += -MMD -MP -MF $@.d
OBJ += $(SRC:%.c=$(OBJDIR)/%.o) $(CXXSRC:%.cpp=$(OBJDIR)/%.o)
DEPENDS = $(SRC:%.c=$(OBJDIR)/%.o.d) $(CXXSRC:%.cpp=$(OBJDIR)/%.o.d)

# Add target to build library
all:lib

lib:$(LIB)

$(LIB):$(OBJ) $(DEPEND_TARGET)
	$(LD) $(LDFLAGS) -o $@ $(OBJ) $(LIBS)

$(OBJDIR)/%.o : %.c
	$(CC) $(CPPFLAGS) $(CFLAGS) -c $< -o $@
$(OBJDIR)/%.o : %.cpp
	$(CXX) $(CPPFLAGS) $(CXXFLAGS) -c $< -o $@
$(OBJDIR)/%.res : %.rc
	$(WINDRES) $< -O coff -o $@ $(RESFLAGS)

clean:
	-rm $(OBJ)
	-rm $(DEPENDS)
	-rmdir $(OBJDIR)
	-rm $(LIB)

.PHONY:clean all lib

-include $(shell mkdir $(OBJDIR) 2>/dev/null) $(wildcard $(OBJDIR)/*.d)
