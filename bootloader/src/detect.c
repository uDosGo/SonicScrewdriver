#include "detect.h"
#include "teletext.h"
#include <string.h>

/* Forward declaration for framebuffer detection */
extern int framebuffer_is_uefi(void);

platform_type_t detect_platform(void) {
    int is_apple = detect_is_apple();
    int has_uefi = detect_has_uefi();

    if (detect_is_vm()) {
        return PLATFORM_VIRTUAL;
    }

    if (is_apple && has_uefi) {
        return PLATFORM_MAC_UEFI;
    }

    if (has_uefi) {
        return PLATFORM_PC_UEFI;
    }

    if (detect_is_bios()) {
        return PLATFORM_PC_BIOS;
    }

    return PLATFORM_UNKNOWN;
}

int detect_is_apple(void) {
    /* Scan SMBIOS for "Apple" in manufacturer string */
    const char *manufacturer = detect_smbios_string(SMBIOS_SYSTEM_INFORMATION, 4);
    if (!manufacturer) return 0;

    /* Check for Apple identifiers */
    if (strstr(manufacturer, "Apple") != NULL) return 1;
    if (strstr(manufacturer, "apple") != NULL) return 1;

    /* Also check product name for Mac identifiers */
    const char *product = detect_smbios_string(SMBIOS_SYSTEM_INFORMATION, 5);
    if (product && strstr(product, "Mac") != NULL) return 1;

    return 0;
}

int detect_is_vm(void) {
    const char *manufacturer = detect_smbios_string(SMBIOS_SYSTEM_INFORMATION, 4);
    if (!manufacturer) return 0;

    if (strstr(manufacturer, "QEMU") != NULL) return 1;
    if (strstr(manufacturer, "VirtualBox") != NULL) return 1;
    if (strstr(manufacturer, "VMware") != NULL) return 1;
    if (strstr(manufacturer, "KVM") != NULL) return 1;
    if (strstr(manufacturer, "Bochs") != NULL) return 1;
    if (strstr(manufacturer, "Xen") != NULL) return 1;
    if (strstr(manufacturer, "Microsoft") != NULL) {
        const char *product = detect_smbios_string(SMBIOS_SYSTEM_INFORMATION, 5);
        if (product && strstr(product, "Virtual") != NULL) return 1;
    }

    return 0;
}

const char *detect_smbios_string(uint8_t type, uint8_t offset) {
    /* Find SMBIOS entry point in BIOS memory area (0xF0000 - 0xFFFFF) */
    /* This is a simplified implementation — real bootloader would scan */
    /* the UEFI configuration tables or BIOS memory for the SMBIOS entry */
    (void)type;
    (void)offset;
    return NULL; /* Will be implemented with actual SMBIOS table walking */
}

uint8_t detect_smbios_version_major(void) {
    return 0;
}

uint8_t detect_smbios_version_minor(void) {
    return 0;
}

int detect_has_uefi(void) {
    /* Check if we're running as a UEFI application */
    /* In UEFI mode, the firmware provides a system table pointer */
    return framebuffer_is_uefi();
}

int detect_is_bios(void) {
    /* If we're not UEFI, we're BIOS */
    return !detect_has_uefi();
}

const char *detect_platform_name(platform_type_t platform) {
    switch (platform) {
        case PLATFORM_MAC_UEFI: return "Mac";
        case PLATFORM_PC_UEFI:  return "PC (UEFI)";
        case PLATFORM_PC_BIOS:  return "PC (BIOS)";
        case PLATFORM_VIRTUAL:  return "Virtual Machine";
        default:                return "Unknown";
    }
}

void detect_debug_dump(void) {
    /* Debug output would go to teletext screen */
    /* Implementation depends on the bootloader's debug output mechanism */
}
