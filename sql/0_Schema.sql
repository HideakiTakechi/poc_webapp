DROP TABLE IF EXISTS `event_image_tagging`;
DROP TABLE IF EXISTS `images`;
DROP TABLE IF EXISTS `event_person_tagging`;
DROP TABLE IF EXISTS `persons`;
DROP TABLE IF EXISTS `events`;
DROP TABLE IF EXISTS `accounts`;

CREATE TABLE `accounts` (
  `account_id` bigint AUTO_INCREMENT,
  `login_name` VARCHAR(255) NOT NULL,
  `shadow_password` VARCHAR(255),
  `mail_address` VARCHAR(255),
  `created_at` DATETIME(6) DEFAULT CURRENT_TIMESTAMP(6),
  `updated_at` DATETIME(6) DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
  PRIMARY KEY(`account_id`)
) ENGINE=InnoDB DEFAULT CHARACTER SET=utf8mb4;

CREATE TABLE `events` (
  `event_id` bigint AUTO_INCREMENT,
  `account_id` bigint NOT NULL,
  `title` VARCHAR(255) NOT NULL,
  `description` TEXT NOT NULL,
  `event_date` DATETIME(6) DEFAULT CURRENT_TIMESTAMP(6),
  PRIMARY KEY(`event_id`),
  INDEX name (`event_date`),
  FOREIGN KEY (`account_id`) REFERENCES accounts(`account_id`) 
) ENGINE=InnoDB DEFAULT CHARACTER SET=utf8mb4;

CREATE TABLE `images` (
  `image_id` bigint AUTO_INCREMENT,
  `image_name` VARCHAR(255) NOT NULL,
  `mime_type` VARCHAR(255) NOT NULL,
  PRIMARY KEY(`image_id`)
);

CREATE TABLE `event_image_tagging` (
  `event_id` bigint NOT NULL,
  `image_id` bigint NOT NULL,
  FOREIGN KEY (`event_id`) REFERENCES events(`event_id`),
  FOREIGN KEY (`image_id`) REFERENCES images(`image_id`),
  INDEX name (`event_id`,`image_id`) );

CREATE TABLE `persons` (
  `person_id` bigint AUTO_INCREMENT,
  `first_name` VARCHAR(255) NOT NULL,
  `last_name` VARCHAR(255) NOT NULL,
  PRIMARY KEY(`person_id`)
);

CREATE TABLE `event_person_tagging` (
  `event_id` bigint NOT NULL,
  `person_id` bigint NOT NULL,
  FOREIGN KEY (`event_id`) REFERENCES events(`event_id`),
  FOREIGN KEY (`person_id`) REFERENCES persons(`person_id`),
  INDEX name (`event_id`,`person_id`)  );
