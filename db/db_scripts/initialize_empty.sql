-- MySQL Workbench Forward Engineering

SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0;
SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0;
SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION';

-- -----------------------------------------------------
-- Schema minitorrent
-- -----------------------------------------------------
DROP SCHEMA IF EXISTS `minitorrent` ;

-- -----------------------------------------------------
-- Schema minitorrent
-- -----------------------------------------------------
CREATE SCHEMA IF NOT EXISTS `minitorrent` ;
USE `minitorrent` ;

-- -----------------------------------------------------
-- Table `minitorrent`.`users`
-- -----------------------------------------------------
DROP TABLE IF EXISTS `minitorrent`.`users` ;

CREATE TABLE IF NOT EXISTS `minitorrent`.`users` (
  `user_ID` INT NOT NULL AUTO_INCREMENT,
  `admin_level` INT NOT NULL COMMENT 'Admin level from 0 - super admin \nlevel  99 - basic user\n',
  `username` VARCHAR(90) NOT NULL,
  `banner_image` VARCHAR(90) NULL DEFAULT 'static/img/default-banner.jpg',
  `profile_pic` VARCHAR(90) NOT NULL DEFAULT '/static/img/default-user.svg',
  `password` VARCHAR(90) NOT NULL,
  `created` DATETIME NOT NULL DEFAULT NOW(),
  `disabled` TINYINT NOT NULL DEFAULT 0,
  `tagline` VARCHAR(45) NULL,
  `bio` TEXT NULL,
  `email` VARCHAR(90) NULL,
  `session_uid` VARCHAR(90) NULL,
  `session_expiry` DATETIME NULL,
  `gender` VARCHAR(45) NULL,
  `session_ip` VARCHAR(45) NULL,
  `user_badges_blob` JSON NULL COMMENT 'denormalized collection of user badges for better performance, these badges are not so important-purely cosmetic and can be updated infrequently to improve database performance',
  PRIMARY KEY (`user_ID`),
  UNIQUE INDEX `username_UNIQUE` (`username` ASC) VISIBLE)
ENGINE = InnoDB;


-- -----------------------------------------------------
-- Table `minitorrent`.`groups`
-- -----------------------------------------------------
DROP TABLE IF EXISTS `minitorrent`.`groups` ;

CREATE TABLE IF NOT EXISTS `minitorrent`.`groups` (
  `group_ID` INT NOT NULL AUTO_INCREMENT,
  `group_name` VARCHAR(90) NOT NULL,
  `group_icon` VARCHAR(90) NULL,
  `group_visibility` VARCHAR(45) CHARACTER SET 'utf8' COLLATE 'utf8_general_ci' NOT NULL,
  `description` TEXT CHARACTER SET 'utf8' NULL COMMENT 'Intro text of the group',
  `join_type` VARCHAR(45) CHARACTER SET 'utf8' COLLATE 'utf8_general_ci' NOT NULL COMMENT 'Must define a value: \"Invite\", \"Invite Request\", \"Free Join\"\n',
  `tagline` TEXT NULL COMMENT 'Short description displayed in lists',
  PRIMARY KEY (`group_ID`))
ENGINE = InnoDB;


-- -----------------------------------------------------
-- Table `minitorrent`.`torrents`
-- -----------------------------------------------------
DROP TABLE IF EXISTS `minitorrent`.`torrents` ;

CREATE TABLE IF NOT EXISTS `minitorrent`.`torrents` (
  `torrent_ID` INT NOT NULL AUTO_INCREMENT,
  `uploaded` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `user_ID` INT NOT NULL,
  `name` TEXT NOT NULL COMMENT 'Name of the folder or the file of the torrent',
  `size` VARCHAR(45) NOT NULL COMMENT 'Human readable size',
  `anonymous` TINYINT NOT NULL DEFAULT 0,
  `access_type` VARCHAR(45) NOT NULL COMMENT 'Selects how the torrent is accessed.\nWhen group ID is set then torrent is visible only inside group\nWhen  field is set to \"WWW\" then torrent is tracked publicly and no authentication or user tracking is performed\nWhen value is set to \"Members\" or any other non \"WWW\" value then only selected site users can access the torrent. For this user IP and unique to user tracking url is cross refrenced.\nWhen access type \"Members Unlisted\" is set then all members can access this torrent through site torrent file download link. User uploads and downloads will be tracked. But torrent will not be listed anywhere but the download page.',
  `group_ID` INT NULL,
  `upvotes` INT NOT NULL DEFAULT 0,
  `downvotes` INT NOT NULL DEFAULT 0,
  `description` TEXT NOT NULL,
  `info_hash` VARCHAR(45) NOT NULL,
  `comment` TEXT NULL,
  `pieces` JSON NULL,
  `piece_length` INT NULL,
  `path` JSON NULL,
  `seeders` INT NOT NULL DEFAULT 0 COMMENT 'aggregate field',
  `leechers` INT NOT NULL DEFAULT 0 COMMENT 'aggregate field',
  `info_field` LONGBLOB NOT NULL,
  `uuid` VARCHAR(40) NOT NULL,
  `category_ID` INT NOT NULL,
  `announce_list` JSON NULL,
  `keep_trackers` TINYINT NOT NULL,
  PRIMARY KEY (`torrent_ID`),
  INDEX `fk_torrents_users1_idx` (`user_ID` ASC) VISIBLE,
  INDEX `fk_torrents_groups1_idx` (`group_ID` ASC) VISIBLE,
  CONSTRAINT `fk_torrents_users1`
    FOREIGN KEY (`user_ID`)
    REFERENCES `minitorrent`.`users` (`user_ID`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION,
  CONSTRAINT `fk_torrents_groups1`
    FOREIGN KEY (`group_ID`)
    REFERENCES `minitorrent`.`groups` (`group_ID`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION)
ENGINE = InnoDB;


-- -----------------------------------------------------
-- Table `minitorrent`.`tags`
-- -----------------------------------------------------
DROP TABLE IF EXISTS `minitorrent`.`tags` ;

CREATE TABLE IF NOT EXISTS `minitorrent`.`tags` (
  `tag_ID` INT NOT NULL AUTO_INCREMENT,
  `tag_name` VARCHAR(45) NOT NULL,
  `color` VARCHAR(45) NULL,
  PRIMARY KEY (`tag_ID`))
ENGINE = InnoDB;


-- -----------------------------------------------------
-- Table `minitorrent`.`torrent_tags`
-- -----------------------------------------------------
DROP TABLE IF EXISTS `minitorrent`.`torrent_tags` ;

CREATE TABLE IF NOT EXISTS `minitorrent`.`torrent_tags` (
  `torrent_tags_ID` INT NOT NULL AUTO_INCREMENT,
  `torrents_torrent_ID` INT NOT NULL,
  `tags_tag_ID` INT NOT NULL,
  INDEX `fk_torrent_tags_tags1_idx` (`tags_tag_ID` ASC) VISIBLE,
  INDEX `fk_torrent_tags_torrents_idx` (`torrents_torrent_ID` ASC) VISIBLE,
  PRIMARY KEY (`torrent_tags_ID`),
  CONSTRAINT `fk_torrent_tags_torrents`
    FOREIGN KEY (`torrents_torrent_ID`)
    REFERENCES `minitorrent`.`torrents` (`torrent_ID`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION,
  CONSTRAINT `fk_torrent_tags_tags1`
    FOREIGN KEY (`tags_tag_ID`)
    REFERENCES `minitorrent`.`tags` (`tag_ID`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION)
ENGINE = InnoDB;


-- -----------------------------------------------------
-- Table `minitorrent`.`torrent_access_list`
-- -----------------------------------------------------
DROP TABLE IF EXISTS `minitorrent`.`torrent_access_list` ;

CREATE TABLE IF NOT EXISTS `minitorrent`.`torrent_access_list` (
  `access_ID` INT NOT NULL AUTO_INCREMENT,
  `user_ID` INT NOT NULL,
  `unique_tracking_url` VARCHAR(45) NOT NULL,
  `torrent_ID` INT NOT NULL,
  `total_uploaded` INT NULL,
  `total_downloaded` INT NULL,
  INDEX `fk_user_access_list_users1_idx` (`user_ID` ASC) VISIBLE,
  PRIMARY KEY (`access_ID`),
  CONSTRAINT `fk_user_access_list_users1`
    FOREIGN KEY (`user_ID`)
    REFERENCES `minitorrent`.`users` (`user_ID`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION,
  CONSTRAINT `fk_user_access_list_torrents1`
    FOREIGN KEY (`torrent_ID`)
    REFERENCES `minitorrent`.`torrents` (`torrent_ID`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION)
ENGINE = InnoDB;


-- -----------------------------------------------------
-- Table `minitorrent`.`friends`
-- -----------------------------------------------------
DROP TABLE IF EXISTS `minitorrent`.`friends` ;

CREATE TABLE IF NOT EXISTS `minitorrent`.`friends` (
  `friends_record_ID` INT NOT NULL AUTO_INCREMENT,
  `user_ID` INT NOT NULL,
  `friend_ID` INT NOT NULL COMMENT 'This table has 2 entries for each friendship link from user > friend and firend > user. This is so that blocking a friend would be directional.',
  `since` DATETIME NOT NULL DEFAULT NOW(),
  `blocked` TINYINT NOT NULL DEFAULT 0 COMMENT '0 - normal friendshipt\n1 - blocked\n',
  `friends_since` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  INDEX `fk_friends_users2_idx` (`friend_ID` ASC) VISIBLE,
  PRIMARY KEY (`friends_record_ID`),
  CONSTRAINT `fk_friends_users1`
    FOREIGN KEY (`user_ID`)
    REFERENCES `minitorrent`.`users` (`user_ID`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION,
  CONSTRAINT `fk_friends_users2`
    FOREIGN KEY (`friend_ID`)
    REFERENCES `minitorrent`.`users` (`user_ID`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION)
ENGINE = InnoDB;


-- -----------------------------------------------------
-- Table `minitorrent`.`group_members`
-- -----------------------------------------------------
DROP TABLE IF EXISTS `minitorrent`.`group_members` ;

CREATE TABLE IF NOT EXISTS `minitorrent`.`group_members` (
  `group_members_ID` INT NOT NULL AUTO_INCREMENT,
  `group_ID` INT NOT NULL,
  `user_ID` INT NOT NULL,
  `group_role` VARCHAR(45) NOT NULL,
  INDEX `fk_group_members_users1_idx` (`user_ID` ASC) VISIBLE,
  PRIMARY KEY (`group_members_ID`),
  CONSTRAINT `fk_group_members_groups1`
    FOREIGN KEY (`group_ID`)
    REFERENCES `minitorrent`.`groups` (`group_ID`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION,
  CONSTRAINT `fk_group_members_users1`
    FOREIGN KEY (`user_ID`)
    REFERENCES `minitorrent`.`users` (`user_ID`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION)
ENGINE = InnoDB;


-- -----------------------------------------------------
-- Table `minitorrent`.`dm_threads`
-- -----------------------------------------------------
DROP TABLE IF EXISTS `minitorrent`.`dm_threads` ;

CREATE TABLE IF NOT EXISTS `minitorrent`.`dm_threads` (
  `dm_thread_ID` INT NOT NULL,
  `thread_title` TEXT NOT NULL COMMENT 'Denormalized field for chat participant list, or a custom name',
  `last_message_date` DATETIME NOT NULL COMMENT 'denormalized field',
  `last_message` TEXT NOT NULL COMMENT 'denormalized field',
  PRIMARY KEY (`dm_thread_ID`))
ENGINE = InnoDB;


-- -----------------------------------------------------
-- Table `minitorrent`.`direct_messages`
-- -----------------------------------------------------
DROP TABLE IF EXISTS `minitorrent`.`direct_messages` ;

CREATE TABLE IF NOT EXISTS `minitorrent`.`direct_messages` (
  `message_ID` INT NOT NULL AUTO_INCREMENT,
  `sender_ID` INT NOT NULL,
  `content` TEXT NOT NULL,
  `date` DATETIME NOT NULL DEFAULT NOW(),
  `dm_thread_ID` INT NOT NULL,
  PRIMARY KEY (`message_ID`),
  INDEX `fk_direct_messages_users1_idx` (`sender_ID` ASC) VISIBLE,
  INDEX `fk_direct_messages_dm_threads1_idx` (`dm_thread_ID` ASC) VISIBLE,
  CONSTRAINT `fk_direct_messages_users1`
    FOREIGN KEY (`sender_ID`)
    REFERENCES `minitorrent`.`users` (`user_ID`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION,
  CONSTRAINT `fk_direct_messages_dm_threads1`
    FOREIGN KEY (`dm_thread_ID`)
    REFERENCES `minitorrent`.`dm_threads` (`dm_thread_ID`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION)
ENGINE = InnoDB;


-- -----------------------------------------------------
-- Table `minitorrent`.`user_torrent_access`
-- -----------------------------------------------------
DROP TABLE IF EXISTS `minitorrent`.`user_torrent_access` ;

CREATE TABLE IF NOT EXISTS `minitorrent`.`user_torrent_access` (
  `session_ID` INT NOT NULL AUTO_INCREMENT,
  `user_ID` INT NOT NULL,
  `uploaded` INT NULL,
  `downloaded` INT NULL,
  `torrent_ID` INT NOT NULL,
  `ip_addr` VARCHAR(45) NOT NULL,
  `start_date` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `end_date` DATETIME NULL,
  `unique_tracking_url` VARCHAR(90) NOT NULL,
  `peer_id` VARCHAR(45) NULL,
  `port` INT NULL,
  PRIMARY KEY (`session_ID`),
  INDEX `fk_user_session_data_users1_idx` (`user_ID` ASC) VISIBLE,
  INDEX `fk_user_session_data_group_torrents1_idx` (`torrent_ID` ASC) VISIBLE,
  UNIQUE INDEX `unique_tracking_url_UNIQUE` (`unique_tracking_url` ASC) VISIBLE,
  CONSTRAINT `fk_user_session_data_users1`
    FOREIGN KEY (`user_ID`)
    REFERENCES `minitorrent`.`users` (`user_ID`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION,
  CONSTRAINT `fk_user_session_data_group_torrents1`
    FOREIGN KEY (`torrent_ID`)
    REFERENCES `minitorrent`.`torrents` (`torrent_ID`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION)
ENGINE = InnoDB;


-- -----------------------------------------------------
-- Table `minitorrent`.`site_news`
-- -----------------------------------------------------
DROP TABLE IF EXISTS `minitorrent`.`site_news` ;

CREATE TABLE IF NOT EXISTS `minitorrent`.`site_news` (
  `post_ID` INT NOT NULL AUTO_INCREMENT,
  `title` TEXT NOT NULL,
  `content` TEXT NOT NULL,
  `posted_by` VARCHAR(45) NOT NULL DEFAULT 'Admin',
  `date` DATETIME NOT NULL DEFAULT NOW(),
  `excerpt` TEXT NOT NULL,
  `commenting` TINYINT NOT NULL DEFAULT 1,
  PRIMARY KEY (`post_ID`))
ENGINE = InnoDB;


-- -----------------------------------------------------
-- Table `minitorrent`.`group_posts`
-- -----------------------------------------------------
DROP TABLE IF EXISTS `minitorrent`.`group_posts` ;

CREATE TABLE IF NOT EXISTS `minitorrent`.`group_posts` (
  `post_ID` INT NOT NULL AUTO_INCREMENT,
  `title` TEXT NOT NULL DEFAULT 'Untitled',
  `content` TEXT NOT NULL DEFAULT ' ',
  `date` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `user_ID` INT NOT NULL,
  `group_ID` INT NOT NULL,
  `group_name` VARCHAR(90) NOT NULL COMMENT 'denormalized field',
  `sticky` TINYINT NOT NULL DEFAULT 0,
  `updated` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `last_reply` DATETIME NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'if null then no replies\n',
  `username` VARCHAR(90) NOT NULL COMMENT 'denormalized field',
  `profile_pic` VARCHAR(90) NOT NULL COMMENT 'denormalized field',
  `reply_count` INT NOT NULL DEFAULT 0 COMMENT 'denormalized field',
  PRIMARY KEY (`post_ID`),
  INDEX `fk_site_news_copy1_users1_idx` (`user_ID` ASC) VISIBLE,
  INDEX `fk_group_posts_groups1_idx` (`group_ID` ASC) VISIBLE,
  CONSTRAINT `fk_site_news_copy1_users1`
    FOREIGN KEY (`user_ID`)
    REFERENCES `minitorrent`.`users` (`user_ID`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION,
  CONSTRAINT `fk_group_posts_groups1`
    FOREIGN KEY (`group_ID`)
    REFERENCES `minitorrent`.`groups` (`group_ID`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION)
ENGINE = InnoDB;


-- -----------------------------------------------------
-- Table `minitorrent`.`user_torrent_ratings`
-- -----------------------------------------------------
DROP TABLE IF EXISTS `minitorrent`.`user_torrent_ratings` ;

CREATE TABLE IF NOT EXISTS `minitorrent`.`user_torrent_ratings` (
  `vote_ID` INT NOT NULL AUTO_INCREMENT,
  `torrents_torrent_ID` INT NOT NULL,
  `users_user_ID` INT NOT NULL,
  `upvote` TINYINT NOT NULL DEFAULT 1,
  PRIMARY KEY (`vote_ID`),
  INDEX `fk_t_upvotes_group_torrents1_idx` (`torrents_torrent_ID` ASC) VISIBLE,
  INDEX `fk_t_upvotes_users1_idx` (`users_user_ID` ASC) VISIBLE,
  CONSTRAINT `fk_t_upvotes_group_torrents1`
    FOREIGN KEY (`torrents_torrent_ID`)
    REFERENCES `minitorrent`.`torrents` (`torrent_ID`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION,
  CONSTRAINT `fk_t_upvotes_users1`
    FOREIGN KEY (`users_user_ID`)
    REFERENCES `minitorrent`.`users` (`user_ID`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION)
ENGINE = InnoDB;


-- -----------------------------------------------------
-- Table `minitorrent`.`badges`
-- -----------------------------------------------------
DROP TABLE IF EXISTS `minitorrent`.`badges` ;

CREATE TABLE IF NOT EXISTS `minitorrent`.`badges` (
  `badge_ID` INT NOT NULL AUTO_INCREMENT,
  `badge_title` VARCHAR(45) NOT NULL,
  `badge_color` VARCHAR(45) NOT NULL,
  PRIMARY KEY (`badge_ID`))
ENGINE = InnoDB;


-- -----------------------------------------------------
-- Table `minitorrent`.`user_badges`
-- -----------------------------------------------------
DROP TABLE IF EXISTS `minitorrent`.`user_badges` ;

CREATE TABLE IF NOT EXISTS `minitorrent`.`user_badges` (
  `user_badges_ID` INT NOT NULL AUTO_INCREMENT,
  `users_user_ID` INT NOT NULL,
  `badges_badge_ID` INT NOT NULL,
  `visible` TINYINT NOT NULL DEFAULT 0,
  PRIMARY KEY (`user_badges_ID`),
  INDEX `fk_user_badges_badges1_idx` (`badges_badge_ID` ASC) VISIBLE,
  CONSTRAINT `fk_user_badges_users1`
    FOREIGN KEY (`users_user_ID`)
    REFERENCES `minitorrent`.`users` (`user_ID`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION,
  CONSTRAINT `fk_user_badges_badges1`
    FOREIGN KEY (`badges_badge_ID`)
    REFERENCES `minitorrent`.`badges` (`badge_ID`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION)
ENGINE = InnoDB;


-- -----------------------------------------------------
-- Table `minitorrent`.`strikes`
-- -----------------------------------------------------
DROP TABLE IF EXISTS `minitorrent`.`strikes` ;

CREATE TABLE IF NOT EXISTS `minitorrent`.`strikes` (
  `strike_ID` INT NOT NULL AUTO_INCREMENT,
  `user_ID` INT NOT NULL,
  `heading` TEXT NOT NULL,
  `description` TEXT NOT NULL,
  `date` DATETIME NOT NULL DEFAULT NOW(),
  PRIMARY KEY (`strike_ID`),
  INDEX `fk_strikes_users1_idx` (`user_ID` ASC) VISIBLE,
  CONSTRAINT `fk_strikes_users1`
    FOREIGN KEY (`user_ID`)
    REFERENCES `minitorrent`.`users` (`user_ID`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION)
ENGINE = InnoDB;


-- -----------------------------------------------------
-- Table `minitorrent`.`comments`
-- -----------------------------------------------------
DROP TABLE IF EXISTS `minitorrent`.`comments` ;

CREATE TABLE IF NOT EXISTS `minitorrent`.`comments` (
  `comment_ID` INT NOT NULL AUTO_INCREMENT,
  `post_ID` INT NOT NULL,
  `user_ID` INT NOT NULL,
  `content` TEXT NOT NULL,
  PRIMARY KEY (`comment_ID`),
  INDEX `fk_comments_site_news1_idx` (`post_ID` ASC) VISIBLE,
  INDEX `fk_comments_users1_idx` (`user_ID` ASC) VISIBLE)
ENGINE = MyISAM;


-- -----------------------------------------------------
-- Table `minitorrent`.`reported_torrents`
-- -----------------------------------------------------
DROP TABLE IF EXISTS `minitorrent`.`reported_torrents` ;

CREATE TABLE IF NOT EXISTS `minitorrent`.`reported_torrents` (
  `report_ID` INT NOT NULL AUTO_INCREMENT,
  `group_torrents_torrent_ID` INT NULL,
  `description` TEXT NOT NULL,
  `reporting_user_ID` INT NOT NULL,
  PRIMARY KEY (`report_ID`),
  INDEX `fk_reported_torrents_users1_idx` (`reporting_user_ID` ASC) VISIBLE,
  CONSTRAINT `fk_reported_torrents_group_torrents1`
    FOREIGN KEY (`group_torrents_torrent_ID`)
    REFERENCES `minitorrent`.`torrents` (`torrent_ID`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION,
  CONSTRAINT `fk_reported_torrents_users1`
    FOREIGN KEY (`reporting_user_ID`)
    REFERENCES `minitorrent`.`users` (`user_ID`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION)
ENGINE = InnoDB;


-- -----------------------------------------------------
-- Table `minitorrent`.`invites`
-- -----------------------------------------------------
DROP TABLE IF EXISTS `minitorrent`.`invites` ;

CREATE TABLE IF NOT EXISTS `minitorrent`.`invites` (
  `invite_ID` INT NOT NULL AUTO_INCREMENT,
  `inviting_user_ID` INT NOT NULL,
  `invite_code` VARCHAR(45) NOT NULL,
  `invited_user_ID` INT NULL,
  `date` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`invite_ID`),
  INDEX `fk_invites_users1_idx` (`inviting_user_ID` ASC) VISIBLE,
  INDEX `fk_invites_users2_idx` (`invited_user_ID` ASC) VISIBLE,
  UNIQUE INDEX `invite_code_UNIQUE` (`invite_code` ASC) VISIBLE,
  CONSTRAINT `fk_invites_users1`
    FOREIGN KEY (`inviting_user_ID`)
    REFERENCES `minitorrent`.`users` (`user_ID`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION,
  CONSTRAINT `fk_invites_users2`
    FOREIGN KEY (`invited_user_ID`)
    REFERENCES `minitorrent`.`users` (`user_ID`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION)
ENGINE = InnoDB;


-- -----------------------------------------------------
-- Table `minitorrent`.`dm_thread_users`
-- -----------------------------------------------------
DROP TABLE IF EXISTS `minitorrent`.`dm_thread_users` ;

CREATE TABLE IF NOT EXISTS `minitorrent`.`dm_thread_users` (
  `entry_ID` INT NOT NULL,
  `user_ID` INT NOT NULL,
  `dm_thread_ID` INT NOT NULL,
  PRIMARY KEY (`entry_ID`),
  INDEX `fk_dm_thread_users_users1_idx` (`user_ID` ASC) VISIBLE,
  INDEX `fk_dm_thread_users_dm_threads1_idx` (`dm_thread_ID` ASC) VISIBLE,
  CONSTRAINT `fk_dm_thread_users_users1`
    FOREIGN KEY (`user_ID`)
    REFERENCES `minitorrent`.`users` (`user_ID`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION,
  CONSTRAINT `fk_dm_thread_users_dm_threads1`
    FOREIGN KEY (`dm_thread_ID`)
    REFERENCES `minitorrent`.`dm_threads` (`dm_thread_ID`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION)
ENGINE = InnoDB;


-- -----------------------------------------------------
-- Table `minitorrent`.`group_post_replies`
-- -----------------------------------------------------
DROP TABLE IF EXISTS `minitorrent`.`group_post_replies` ;

CREATE TABLE IF NOT EXISTS `minitorrent`.`group_post_replies` (
  `reply_ID` INT NOT NULL,
  `post_ID` INT NOT NULL,
  `user_ID` INT NOT NULL,
  `content` TEXT NOT NULL,
  `date` DATETIME NOT NULL,
  `updated` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`reply_ID`),
  INDEX `fk_group_post_replies_group_posts1_idx` (`post_ID` ASC) VISIBLE,
  INDEX `fk_group_post_replies_users1_idx` (`user_ID` ASC) VISIBLE,
  CONSTRAINT `fk_group_post_replies_group_posts1`
    FOREIGN KEY (`post_ID`)
    REFERENCES `minitorrent`.`group_posts` (`post_ID`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION,
  CONSTRAINT `fk_group_post_replies_users1`
    FOREIGN KEY (`user_ID`)
    REFERENCES `minitorrent`.`users` (`user_ID`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION)
ENGINE = InnoDB;


-- -----------------------------------------------------
-- Table `minitorrent`.`torrent_comments`
-- -----------------------------------------------------
DROP TABLE IF EXISTS `minitorrent`.`torrent_comments` ;

CREATE TABLE IF NOT EXISTS `minitorrent`.`torrent_comments` (
  `comment_ID` INT NOT NULL AUTO_INCREMENT,
  `torrent_ID` INT NOT NULL,
  `user_ID` INT NOT NULL,
  `content` TEXT NOT NULL,
  `date` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`comment_ID`),
  INDEX `fk_torrent_comments_users1_idx` (`user_ID` ASC) VISIBLE,
  CONSTRAINT `fk_torrent_comments_torrents1`
    FOREIGN KEY (`torrent_ID`)
    REFERENCES `minitorrent`.`torrents` (`torrent_ID`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION,
  CONSTRAINT `fk_torrent_comments_users1`
    FOREIGN KEY (`user_ID`)
    REFERENCES `minitorrent`.`users` (`user_ID`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION)
ENGINE = InnoDB;


-- -----------------------------------------------------
-- Table `minitorrent`.`friend_requests`
-- -----------------------------------------------------
DROP TABLE IF EXISTS `minitorrent`.`friend_requests` ;

CREATE TABLE IF NOT EXISTS `minitorrent`.`friend_requests` (
  `friend_request_ID` INT NOT NULL AUTO_INCREMENT,
  `receiver_user_ID` INT NOT NULL,
  `sender_user_ID` INT NOT NULL,
  `status` TINYINT NOT NULL COMMENT '0 - pending\n1 - declined',
  `message` TEXT NULL,
  `date` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`friend_request_ID`),
  INDEX `fk_friend_requests_users1_idx` (`receiver_user_ID` ASC) VISIBLE,
  INDEX `fk_friend_requests_users2_idx` (`sender_user_ID` ASC) VISIBLE,
  CONSTRAINT `fk_friend_requests_users1`
    FOREIGN KEY (`receiver_user_ID`)
    REFERENCES `minitorrent`.`users` (`user_ID`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION,
  CONSTRAINT `fk_friend_requests_users2`
    FOREIGN KEY (`sender_user_ID`)
    REFERENCES `minitorrent`.`users` (`user_ID`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION)
ENGINE = InnoDB;


-- -----------------------------------------------------
-- Table `minitorrent`.`peers`
-- -----------------------------------------------------
DROP TABLE IF EXISTS `minitorrent`.`peers` ;

CREATE TABLE IF NOT EXISTS `minitorrent`.`peers` (
  `entry_ID` INT NOT NULL AUTO_INCREMENT,
  `torrent_ID` INT NOT NULL,
  `peer_id` VARCHAR(45) NOT NULL,
  `ip` VARCHAR(45) NOT NULL,
  `port` INT NOT NULL,
  `last_update` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `left` VARCHAR(140) NOT NULL,
  PRIMARY KEY (`entry_ID`),
  INDEX `fk_peers_torrents1_idx` (`torrent_ID` ASC) VISIBLE,
  UNIQUE INDEX `peer_id_UNIQUE` (`peer_id` ASC) VISIBLE,
  CONSTRAINT `fk_peers_torrents1`
    FOREIGN KEY (`torrent_ID`)
    REFERENCES `minitorrent`.`torrents` (`torrent_ID`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION)
ENGINE = InnoDB
COMMENT = 'last update keeps track when the peer last sent a request to server, delete peers after few hours to avoid sending invalid addresses to other peers.';


-- -----------------------------------------------------
-- Table `minitorrent`.`sys_info`
-- -----------------------------------------------------
DROP TABLE IF EXISTS `minitorrent`.`sys_info` ;

CREATE TABLE IF NOT EXISTS `minitorrent`.`sys_info` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `schema_revision` VARCHAR(45) NOT NULL DEFAULT '0.14',
  `site_name` VARCHAR(145) NOT NULL DEFAULT 'Minitorrent',
  PRIMARY KEY (`id`))
ENGINE = InnoDB;

-- -----------------------------------------------------
-- Data for table `minitorrent`.`sys_info` insert revision data
-- -----------------------------------------------------
START TRANSACTION;
USE `minitorrent`;
INSERT INTO `minitorrent`.`sys_info` (`id`, `schema_revision`, `site_name`) VALUES (1, '0.14', 'Minitorrent');

COMMIT;


SET SQL_MODE=@OLD_SQL_MODE;
SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS;
SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS;

