#ifndef DETECT_H
#define DETECT_H

#include <stdint.h>
#include "../include/menu.h"  /* For platform_type_t */

/* SMBIOS structure types */
#define SMBIOS_BIOS_INFORMATION         0
#define SMBIOS_SYSTEM_INFORMATION       1
#define SMBIOS_BASEBOARD_INFORMATION    2
#define SMBIOS_SYSTEM_ENCLOSURE         3

/* SMBIOS entry point structure */
typedef struct {
    char     anchor[4];          /* "_SM_" */
    uint8_t  checksum;
    uint8_t  length;
    uint8_t  major_version;
    uint8_t  minor_version;
    uint16_t max_structure_size;
    uint8_t  revision;
    uint8_t  formatted_area[5];
    char     intermediate_anchor[5]; /* "_DMI_" */
    uint8_t  intermediate_checksum;
    uint16_t structure_table_length;
    uint32_t structure_table_address;
    uint16_t number_of_structures;
    uint8_t  bcd_revision;
} __attribute__((packed)) smbios_entry_t;

/* SMBIOS header (common to all structures) */
typedef struct {
    uint8_t  type;
    uint8_t  length;
    uint16_t handle;
} __attribute__((packed)) smbios_header_t;

/* SMBIOS System Information (type 1) */
typedef struct {
    smbios_header_t header;
    uint8_t  manufacturer;
    uint8_t  product_name;
    uint8_t  version;
    uint8_t  serial_number;
    uint8_t  uuid[16];
    uint8_t  wake_up_type;
    uint8_t  sku_number;
    uint8_t  family;
} __attribute__((packed)) smbios_system_info_t;

/* ACPI RSDP descriptor */
typedef struct {
    char     signature[8];       /* "RSD PTR " */
    uint8_t  checksum;
    char     oem_id[6];
    uint8_t  revision;
    uint32_t rsdt_address;
} __attribute__((packed)) acpi_rsdp_t;

/* ACPI RSDT (v1) */
typedef struct {
    char     signature[4];       /* "RSDT" */
    uint32_t length;
    uint8_t  revision;
    uint8_t  checksum;
    char     oem_id[6];
    char     oem_table_id[8];
    uint32_t oem_revision;
    uint32_t creator_id;
    uint32_t creator_revision;
} __attribute__((packed)) acpi_sdt_header_t;

/* Function prototypes */

/* Detect platform by scanning SMBIOS and ACPI */
platform_type_t detect_platform(void);

/* Check if running on Apple hardware via SMBIOS */
int detect_is_apple(void);

/* Check if running in a virtual machine */
int detect_is_vm(void);

/* Get SMBIOS system information string (returns NULL if not found) */
const char *detect_smbios_string(uint8_t type, uint8_t offset);

/* Get the SMBIOS major version */
uint8_t detect_smbios_version_major(void);

/* Get the SMBIOS minor version */
uint8_t detect_smbios_version_minor(void);

/* Check if UEFI is available */
int detect_has_uefi(void);

/* Check if running on legacy BIOS */
int detect_is_bios(void);

/* Get a human-readable platform name */
const char *detect_platform_name(platform_type_t platform);

/* Debug: dump SMBIOS info to teletext screen */
void detect_debug_dump(void);

/* -- Lifecycle state reporting (SonicScrewdriver v2.1) -- */

/* Lifecycle state enum aligned to bootloader-status-taxonomy.md */
typedef enum {
    LIFECYCLE_INIT,
    LIFECYCLE_DETECT,
    LIFECYCLE_DETECT_COMPLETE,
    LIFECYCLE_FRAMEBUFFER_INIT,
    LIFECYCLE_FB_FALLBACK,
    LIFECYCLE_CONFIG_LOAD,
    LIFECYCLE_CONFIG_EMPTY,
    LIFECYCLE_THEME_LOAD,
    LIFECYCLE_THEME_DEFAULT,
    LIFECYCLE_MENU_RENDER,
    LIFECYCLE_MENU_TIMEOUT,
    LIFECYCLE_CHAINLOAD_BEGIN,
    LIFECYCLE_CHAINLOAD_FAIL,
    LIFECYCLE_CHAINLOAD_SUCCESS,
    LIFECYCLE_PANIC_OOB,
    LIFECYCLE_PANIC_CORRUPT,
    LIFECYCLE_HALT,
} lifecycle_state_t;

/* Report a lifecycle state transition with platform context.
 * Called at each lifecycle stage to emit structured status
 * suitable for spool event consumption (via serial or EFI var).
 */
void detect_report_lifecycle(lifecycle_state_t state);

/* Get a short name for a lifecycle state (e.g. "detect_complete") */
const char *detect_lifecycle_name(lifecycle_state_t state);

/* Get the severity string: INFO, WARNING, ERROR, CRITICAL */
const char *detect_lifecycle_severity(lifecycle_state_t state);

/* Check if running on ARM64 (device-tree based detection) */
int detect_is_arm64(void);

/* Check if running on Apple Silicon (ARM64 + Apple device-tree) */
int detect_is_apple_silicon(void);

#endif /* DETECT_H */
