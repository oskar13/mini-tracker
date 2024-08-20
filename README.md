# Mini Tracker
A small/mid scale torrent tracker and a site written in Go. It features public and private torrent tracking mode. Web interface gives access to community features like groups and direct messages.

- Invite based account creation
- Simple chat system
- Site news with commenting
- User profiles
- Community hubs
- Public/private torrent tracking and listing

⭐ Leave a star to show support for this project!


## Current progress
> [!NOTE]  
> This project is not fully functional yet and is still in its early development stage.

## How to run

### With docker

Run `gen_secrets.sh` which will generate password files with openssl. Build and run docker containers with `docker-compose up --build `. Navigate to port 8080 on the server to start the web install process (see bellow).

### Locally

#### Requirements:
 * A recent MariaDB (tested with 11.5.2) or MySQL database install.
 * Go >= 1.23

Run `gen_secrets.sh` which will generate password files with openssl. Then edit start.sh with your MySQL credentials (make a non root user with access to a database). Edit `db_password.txt` with the user's password. `db_root_password.txt` is only needed when using docker so this file can be ignored. Run `start.sh` to start the server with go run. Navigate to port 8080 and proceed with the installer (see bellow).

## Web Installer
If no previous data in database is detected then the installer will be launched in website root. When entering the installer, a token is needed to proceed. You can find it in `installer_token.txt` in project root directory after running `gen_secrets.sh`. After successfully validating the token, you can create admin user account, set site name and finish the install.

### MySQL Workbench project
With this repository also comes a MySQL Workbench project with ERD diagram with embedded data to test out site functionality.

## Goals
- Make it easy to share torrent files between small number of people
- Provide a simple to use interface for uploading torrents to a tracker
- Feature a small community hub related to torrents
- User managed access rules for torrent tracking
- Deploy tracker via docker
- Somehow make it easy to use via Discord
- Keep it simple/minimal - avoid using too many dependencies or frameworks
- Keep the site functional even on browsers with no-script extension installed

### What are not the goals of this project
- Promote piracy
- Promote sharing illegal content

## Access to torrents
Torrents can be categorized by their access type.

All public torrents and the list of peers are accessible through the public tracker interface without having an account.

- **Public Listed** torrents will appear in searches and category pages.
- **Public Unlisted** torrents will be accessible through their URL which opens the torrent info page. 

Private torrents require an account to access the info page. To access peers on the tracker, unique tracking URL (generated during torrent file download from the info page) is used to authenticate users. While this approach has some serious security flaws there are no other alternatives with the current bittorrent technology. Main goal here is to limit the information visibility on the web and enable tracking user upload/download ratios.

- **Members Listed** torrents will appear in searches and category pages for site members.
- **Members Unlisted** info page can be accessed by site members who have the URL.
- **Members Access List** info page can be accessed by site members who have the URL and are on the access list set by the uploader.
- **Group Public** if a group is set as Public (any site member can access) then by default, torrents of that group are set as Group Public which means that any site member can view and download them.
- **Group Private** only members of the group can see the torrent.




## Why
Sometimes there is a need to move large files (over 20G - HDD backup images, neural networks, etc.) from one computer to another over the internet with a one-time transaction. In this case, it wouldn't make sense to buy storage space with an expensive cloud storage provider or set up FTP server which requires additional configuration. 

Additionally, while transferring large files, the network connection might become unstable and cause the upload to fail.

Instead one should use a P2P protocol like Bit-torrent which is resilient to potential connection issues.

The problem is that in the current situation, it might be hard to use widely known public trackers to share files between friends or colleagues. Some domains may be blocked by ISPs or draw unwanted attention when connecting to them. Also, there might be a privacy/legal concerns when sharing a file on a public tracker.

The solution is to run your own tracker that is easily deployable when needed. One which features different privacy levels.

---

> [!IMPORTANT]  
> Bit-torrent protocol does not feature any strong privacy features like good encryption or ability to efficiently prevent third parties from downloading the torrent content. Much of this is due the limitations of the protocol and also the bit-torrent client programs themselves. It is presumed that files with sensitive content will be encrypted beforehand by the users. Community feature tied to the tracker can be used to add external "invisible" data to go with the torrent when it is shared between members: like decryption keys, instructions, etc.