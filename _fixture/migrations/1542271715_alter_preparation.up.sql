SET FOREIGN_KEY_CHECKS = 0;
ALTER TABLE `preparation_document`
ADD COLUMN `year` VARCHAR(20) NOT NULL  COMMENT 'batch code with year code only' AFTER `quantity`,
ADD COLUMN `week` VARCHAR(20) NOT NULL  COMMENT 'batch code with week code only' AFTER `quantity`;

ALTER TABLE `preparation_actual`
ADD COLUMN `year` VARCHAR(20) NOT NULL  COMMENT 'batch code with year code only' AFTER `note`,
ADD COLUMN `week` VARCHAR(20) NOT NULL  COMMENT 'batch code with week code only' AFTER `note`;