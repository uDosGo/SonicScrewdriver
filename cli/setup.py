#!/usr/bin/env python3
"""SonicScrewdriver v2 — Universal USB Bootloader & System Toolkit CLI."""

from setuptools import setup, find_packages

setup(
    name="sonic-screwdriver",
    version="2.1.0",
    description="Universal USB Bootloader & System Toolkit",
    author="uDosGo",
    url="https://github.com/uDosGo/SonicScrewdriver",
    packages=find_packages(),
    include_package_data=True,
    install_requires=[
        "click>=8.0",
        "pyyaml>=6.0",
        "pydantic>=2.0",
        "rich>=13.0",
        "requests>=2.31",
        "cryptography>=41.0",
        "fido2>=1.1",
        "python-gnupg>=0.5",
    ],
    extras_require={
        "mesh": ["pynacl>=1.5", "zeroconf>=0.131"],
        "flash": ["pyserial>=3.5", "esptool>=4.0"],
        "dev": ["pytest>=7.0", "black>=23.0", "mypy>=1.0"],
    },
    entry_points={
        "console_scripts": [
            "sonic=sonic.cli:cli",
        ],
    },
    python_requires=">=3.11",
)
