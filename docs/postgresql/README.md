# postgresql table descriptions


Name | Value | Description
--- | --- | ---
id                           | uuid | Domain ID
created                      | timestamp with time zone | Timestamp as domain was created
modified                     | timestamp with time zone | Timestamp as domain was modified
status                       | character varying(20) | Domain status
verbose_name                 | character varying(255) | Domain verbose name
description                  | text | Domain description
memory_count                 | integer | RAM
power_state                  | character varying(20) | D
suspended                    | boolean | VM in suspend state?
controllers                  | jsonb | Availables controllers
video                        | jsonb | fbuf params
boot                         | jsonb | Boot order[1]
storage_pool_id              | uuid | Pool ID
node_id                      | uuid | Node ID
locked_by_id                 | uuid | Locked ID
os_type                      | character varying(10) | Running OS
os_profile                   | character varying(40) | OS profile
boot_type                    | character varying(10) | Boot type
graphics_password            | character varying(16) | Password for graphics
remote_access                | boolean | True if remote access available
remote_access_port           | integer | Port number
remote_access_allow_all      | boolean | True if access allowed
remote_access_white_list     | jsonb | Allowed IP Addresses for VNC
cpu_topology                 | jsonb | CPU topology[2]
pci_devices                  | jsonb | Devices attached to domain[3]


* [1]:
Boot order json sample:


    {
        "a981bf0a-f6e5-4cbc-ba3a-df68ec8c916f": 1,
        "ef4153b7-7333-4f8c-a1fd-766b5b384408": 2
    }


* [2]:

CPU topology json sample:

    {
        "cpu_cores": 1,
        "cpu_sockets": 4,
        "cpu_threads": 1
    }


* [3]:

PCI devices json sample:

    "pci": [
        {
            "model": "hostbridge",
            "bus": 1
        }
    ],
    "usb": [
        {
            "model": "xhci",
            "bus": 31
        }
    ],
    "virtio": [
        {}
    ],
    "scsi": []

