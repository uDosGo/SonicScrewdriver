# SonicScrewdriver v2 — Bootloader Makefile
#
# Targets:
#   make ARCH=x86_64          — Build UEFI bootloader (sonic-bootloader.efi)
#   make ARCH=x86_64 BIOS=1   — Build legacy BIOS bootloader (sonic-bootloader.bin)
#   make qemu-test             — Test in QEMU
#   make clean                 — Clean build artifacts

ARCH ?= x86_64
BIOS ?= 0

# Output directory
BUILD_DIR := build
SRC_DIR := src
INC_DIR := include
CFG_DIR := config

# Toolchain
ifeq ($(ARCH), x86_64)
    CC := gcc
    LD := ld
    OBJCOPY := objcopy
    ARCH_CFLAGS := -m64 -mno-red-zone -mno-mmx -mno-sse -mno-sse2
    ARCH_LDFLAGS := -m elf_x86_64
else ifeq ($(ARCH), aarch64)
    CC := aarch64-linux-gnu-gcc
    LD := aarch64-linux-gnu-ld
    OBJCOPY := aarch64-linux-gnu-objcopy
    ARCH_CFLAGS :=
    ARCH_LDFLAGS :=
else
    $(error Unsupported architecture: $(ARCH))
endif

# Common flags
CFLAGS := $(ARCH_CFLAGS) -ffreestanding -fno-stack-protector -fno-pic -mno-red-zone
CFLAGS += -I$(INC_DIR) -Wall -Wextra -Werror -O2 -std=c11
LDFLAGS := $(ARCH_LDFLAGS) -nostdlib -z noexecstack

# Source files
SRCS := $(wildcard $(SRC_DIR)/*.c)
OBJS := $(SRCS:$(SRC_DIR)/%.c=$(BUILD_DIR)/%.o)

# UEFI-specific
ifeq ($(BIOS), 0)
    CFLAGS += -DUEFI_BUILD -I/usr/include/efi -I/usr/include/efi/$(ARCH)
    LDFLAGS += -T /usr/lib/efi/elf_$(ARCH)_efi.lds
    LIBS := -lgnuefi -lefi
    TARGET := $(BUILD_DIR)/sonic-bootloader.efi
    FINAL_TARGET := sonic-bootloader.efi
else
    CFLAGS += -DBIOS_BUILD -fno-toplevel-reorder
    LDFLAGS += -Ttext 0x7C00 --oformat binary
    TARGET := $(BUILD_DIR)/sonic-bootloader.bin
    FINAL_TARGET := sonic-bootloader.bin
endif

.PHONY: all clean qemu-test dirs

all: dirs $(TARGET)

# Create build directory
dirs:
	@mkdir -p $(BUILD_DIR)

# Compile C sources
$(BUILD_DIR)/%.o: $(SRC_DIR)/%.c $(INC_DIR)/*.h
	$(CC) $(CFLAGS) -c $< -o $@

# Link UEFI binary
ifeq ($(BIOS), 0)
$(TARGET): $(OBJS)
	$(LD) $(LDFLAGS) $^ -o $(BUILD_DIR)/sonic-bootloader.so $(LIBS)
	$(OBJCOPY) -j .text -j .sdata -j .data -j .dynamic -j .dynsym \
	           -j .rel -j .rela -j .reloc --target=efi-app-$(ARCH) \
	           $(BUILD_DIR)/sonic-bootloader.so $@
	@echo "  ✅ Built: $(FINAL_TARGET)"
	@ls -lh $@
else
# Link BIOS binary
$(TARGET): $(OBJS)
	$(LD) $(LDFLAGS) $^ -o $@
	@echo "  ✅ Built: $(FINAL_TARGET)"
	@ls -lh $@
endif

# Test in QEMU
qemu-test: all
ifeq ($(BIOS), 0)
	@echo "  🧪 Testing UEFI bootloader in QEMU..."
	@echo "     Copy $(FINAL_TARGET) to a FAT32 image and boot with:"
	@echo "     qemu-system-x86_64 -bios /usr/share/ovmf/OVMF.fd \\"
	@echo "       -drive file=fat:rw:$(BUILD_DIR)/esp,format=raw"
else
	@echo "  🧪 Testing BIOS bootloader in QEMU..."
	@echo "     qemu-system-x86_64 -drive file=$(TARGET),format=raw"
endif

# Clean
clean:
	rm -rf $(BUILD_DIR)
	@echo "  🧹 Cleaned build artifacts"

# Install prerequisites (macOS)
install-deps-macos:
	@echo "  📦 Installing bootloader build dependencies..."
	brew install gnu-efi mtools xorriso qemu
	@echo "  ✅ Dependencies installed"

# Install prerequisites (Linux)
install-deps-linux:
	@echo "  📦 Installing bootloader build dependencies..."
	sudo apt-get install -y gnu-efi mtools xorriso qemu-system-x86 ovmf
	@echo "  ✅ Dependencies installed"

# Generate embedded font from Linux console font
generate-font:
	@echo "  🔤 Generating embedded font..."
	@echo "     (Requires Linux with console-setup package)"
	@echo "     cp /usr/share/consolefonts/Lat2-Terminus16.psf font.psf"
	@echo "     Then convert to C array"

# Show build configuration
info:
	@echo "SonicScrewdriver Bootloader Build Configuration"
	@echo "================================================"
	@echo "  Architecture:    $(ARCH)"
	@echo "  Target:          $(if $(BIOS),BIOS,UEFI)"
	@echo "  CC:              $(CC)"
	@echo "  CFLAGS:          $(CFLAGS)"
	@echo "  LDFLAGS:         $(LDFLAGS)"
	@echo "  Source files:    $(SRCS)"
	@echo "  Output:          $(FINAL_TARGET)"
	@echo "================================================"
