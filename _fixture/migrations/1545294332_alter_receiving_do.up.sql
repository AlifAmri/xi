SET FOREIGN_KEY_CHECKS = 0;
ALTER TABLE `incoming_vehicle`
CHANGE COLUMN `status` `status` ENUM('in_queue', 'in_progress', 'finished', 'out', 'cancelled')  NOT NULL DEFAULT 'in_queue';

ALTER TABLE `delivery_order`
ADD COLUMN  `is_active` TINYINT(1) NOT NULL DEFAULT 1;

ALTER TABLE `receiving`
ADD COLUMN  `is_active` TINYINT(1) NOT NULL DEFAULT 1;