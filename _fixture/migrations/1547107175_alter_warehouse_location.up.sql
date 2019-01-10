SET FOREIGN_KEY_CHECKS = 0;
ALTER TABLE `warehouse_location`
CHANGE COLUMN `storage_used` `storage_used` INT(11)  NOT NULL DEFAULT 0;
