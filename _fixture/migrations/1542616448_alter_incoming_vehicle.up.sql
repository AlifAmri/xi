SET FOREIGN_KEY_CHECKS = 0;
ALTER TABLE `incoming_vehicle`
CHANGE COLUMN `status` `status` ENUM('in_queue', 'in_progress', 'finished', 'out')  NOT NULL DEFAULT 'in_queue';