# Mini Tracker
A small/mid scale torrent tracker and a site written in Go. It features public and private torrent tracking mode. Web interface gives access to community features like groups and direct messages.

- Invite based account creation
- Simple chat system
- Site news with commenting
- User profiles
- Community hubs
- Public/private torrent tracking and listing
- Developed for Tor browser/Onion Service in mind as a possible deployment platform

⭐ Leave a star to show support for this project!


## Current progress
![60%](https://progress-bar.dev/60)

> [!NOTE]  
> This project is not fully functional yet and is still in its early development stage.

## How to run
Installer script coming soon. To run the code at current state, edit start.sh with your MySQL credentials, make a minitorrent database on your server. 

(Optional) In MySQL workbench forward engineer a schema with insert commands. Log in with user accounts provided in tests/testusers.txt file.


### Goals
- Make it easy to share torrent files between small number of people
- Provide a simple to use interface for uploading torrents to a tracker
- Feature a small community hub related to torrents
- Private torrents accessed through REST like API
- User managed access rules for torrent tracking
- Deploy tracker via docker
- Somehow make it easy to use via Discord
- Keep it simple/minimal - avoid using too many dependencies or frameworks
- Keep the site functional even on browsers with no-script extension installed

### What are not the goals of this project
- Promote piracy
- Promote sharing illegal content

### MySQL Workbench project
With this repository also comes a MySQL Workbench project with ERD diagram with embedded data to test out site functionality.

## Why
Sometimes there is a need to move large files (over 20G - HDD backup images, neural networks, etc.) from one computer to another over the internet with a one-time transaction. In this case, it wouldn't make sense to buy storage space with an expensive cloud storage provider or set up FTP server which requires additional configuration. 

Additionally, while transferring large files, the network connection might become unstable and cause the upload to fail.

Instead one should use a P2P protocol like Bit-torrent which is resilient to potential connection issues.

The problem is that in the current situation, it might be hard to use widely known public trackers to share files between friends or colleagues. Some domains may be blocked by ISPs or draw unwanted attention when connecting to them. Also, there might be a privacy/legal concerns when sharing a file on a public tracker.

The solution is to run your own tracker that is easily deployable when needed. One which features different privacy levels.

---

> [!IMPORTANT]  
> Bit-torrent protocol does not feature any strong privacy features like good encryption or ability to efficiently prevent third parties from downloading the torrent content. Much of this is due the limitations of the protocol and also the bit-torrent client programs themselves. It is presumed that files with sensitive content will be encrypted beforehand by the users. Community feature tied to the tracker can be used to add external "invisible" data to go with the torrent when it is shared between members: like decryption keys, instructions, etc.