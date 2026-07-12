#include "detect.h"
#include "teletext.h"
#include <string.h>

/* Forward declaration for framebuffer detection */
extern int framebuffer_is_uefi(void);

platform_type_t detect_platform(void) {
    detect_report_lifecycle(LIFECYCLE_DETECT);

    int is_apple = detect_is_apple();
    int has_uefi = detect_has_uefi();
    int is_arm = detect_is_arm64();
    int is_vm = detect_is_vm();

    if (is_vm) {
        detect_report_lifecycle(LIFECYCLE_DETECT_COMPLETE);
        return PLATFORM_VIRTUAL;
    }

    if (is_arm) {
        /* ARM64 path — Apple Silicon, Raspberry Pi, etc. */
        detect_report_lifecycle(LIFECYCLE_DETECT_COMPLETE);
        return PLATFORM_ARM64;
    }

    if (is_apple && has_uefi) {
        detect_report_lifecycle(LIFECYCLE_DETECT_COMPLETE);
        return PLATFORM_MAC_UEFI;
    }

    if (has_uefi) {
        detect_report_lifecycle(LIFECYCLE_DETECT_COMPLETE);
        return PLATFORM_PC_UEFI;
    }

    if (detect_is_bios()) {
        detect_report_lifecycle(LIFECYCLE_DETECT_COMPLETE);
        return PLATFORM_PC_BIOS;
    }

    detect_report_lifecycle(LIFECYCLE_DETECT_COMPLETE);
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
        case PLATFORM_MAC_UEFI:    return "Mac";
        case PLATFORM_PC_UEFI:     return "PC (UEFI)";
        case PLATFORM_PC_BIOS:     return "PC (BIOS)";
        case PLATFORM_VIRTUAL:     return "Virtual Machine";
        case PLATFORM_ARM64:       return "ARM64";
        default:                   return "Unknown";
    }
}

void detect_debug_dump(void) {
    /* Debug output would go to teletext screen */
    /* Implementation depends on the bootloader's debug output mechanism */
}

/* -- Lifecycle state reporting (SonicScrewdriver v2.1) -- */

void detect_report_lifecycle(lifecycle_state_t state) {
    /* Writes structured lifecycle event to serial/EFI variable.
     * Format: "SONIC_LIFECYCLE:<name>:<severity>:<platform>"
     * This can be consumed by a host-side spool adapter.
     */
    const char *name = detect_lifecycle_name(state);
    const char *severity = detect_lifecycle_severity(state);
    const char *plat = detect_platform_name(detect_platform());
    (void)name;
    (void)severity;
    (void)plat;
    /* TODO: Write to EFI variable or serial console:
     *   serial_printf("SONIC_LIFECYCLE:%s:%s:%s\n", name, severity, plat);
     */
}

const char *detect_lifecycle_name(lifecycle_state_t state) {
    switch (state) {
        case LIFECYCLE_INIT:             return "init";
        case LIFECYCLE_DETECT:           return "detect";
        case LIFECYCLE_DETECT_COMPLETE:  return "detect_complete";
        case LIFECYCLE_FRAMEBUFFER_INIT: return "framebuffer_init";
        case LIFECYCLE_FB_FALLBACK:      return "fb_fallback";
        case LIFECYCLE_CONFIG_LOAD:      return "config_load";
        case LIFECYCLE_CONFIG_EMPTY:     return "config_empty";
        case LIFECYCLE_THEME_LOAD:       return "theme_load";
        case LIFECYCLE_THEME_DEFAULT:    return "theme_default";
        case LIFECYCLE_MENU_RENDER:      return "menu_render";
        case LIFECYCLE_MENU_TIMEOUT:     return "menu_timeout";
        case LIFECYCLE_CHAINLOAD_BEGIN:  return "chainload_begin";
        case LIFECYCLE_CHAINLOAD_FAIL:   return "chainload_fail";
        case LIFECYCLE_CHAINLOAD_SUCCESS:return "chainload_success";
        case LIFECYCLE_PANIC_OOB:        return "panic_oob";
        case LIFECYCLE_PANIC_CORRUPT:    return "panic_corrupt";
        case LIFECYCLE_HALT:             return "halt";
        default:                         return "unknown";
    }
}

const char *detect_lifecycle_severity(lifecycle_state_t state) {
    switch (state) {
        case LIFECYCLE_INIT:
        case LIFECYCLE_DETECT:
        case LIFECYCLE_DETECT_COMPLETE:
        case LIFECYCLE_FRAMEBUFFER_INIT:
        case LIFECYCLE_CONFIG_LOAD:
        case LIFECYCLE_THEME_DEFAULT:
        case LIFECYCLE_MENU_RENDER:
        case LIFECYCLE_MENU_TIMEOUT:
        case LIFECYCLE_CHAINLOAD_BEGIN:
        case LIFECYCLE_CHAINLOAD_SUCCESS:
            return "INFO";
        case LIFECYCLE_FB_FALLBACK:
        case LIFECYCLE_CONFIG_EMPTY:
            return "WARNING";
        case LIFECYCLE_THEME_LOAD:
        case LIFECYCLE_CHAINLOAD_FAIL:
        case LIFECYCLE_HALT:
            return "ERROR";
        case LIFECYCLE_PANIC_OOB:
        case LIFECYCLE_PANIC_CORRUPT:
            return "CRITICAL";
        default:
            return "INFO";
    }
}

int detect_is_arm64(void) {
    /* Device-tree based ARM64 detection.
     * On ARM64 systems (Apple Silicon, Raspberry Pi), SMBIOS may
     * not be available. Instead, check for device-tree or
     * ACPI ARM-specific tables (GTDT, SPCR).
     *
     * Heuristic: check CPU ID registers or device-tree presence.
     * For now, check if we're NOT x86 (UEFI on ARM64 uses
     * a different machine type in the PE header).
     */
    /* Simplified: if SMBIOS returns nothing and we have UEFI,
     * we may be on ARM64. A proper check would read the
     * device-tree /proc/device-tree or check the UEFI system
     * table's firmware vendor for ARM-specific strings.
     */
    const char *manufacturer = detect_smbios_string(
        SMBIOS_SYSTEM_INFORMATION, 4);
    if (!manufacturer && detect_has_uefi()) {
        /* No SMBIOS but UEFI present — likely ARM64 */
        return 1;
    }

    /* Check for ARM64-specific SMBIOS vendor strings */
    if (manufacturer) {
        if (strstr(manufacturer, "Raspberry") != NULL) return 1;
        if (strstr(manufacturer, "ARM") != NULL) return 1;
    }

    return 0;
}

int detect_is_apple_silicon(void) {
    /* Apple Silicon = ARM64 + Apple manufacturer */
    if (!detect_is_arm64()) return 0;
    return detect_is_apple();
}
