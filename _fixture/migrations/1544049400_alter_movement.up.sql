SET FOREIGN_KEY_CHECKS = 0;
ALTER TABLE `stock_movement`
ADD COLUMN  `pallet_id` BIGINT(11) UNSIGNED NULL,
ADD INDEX `movement_ibfk_8_idx` (`pallet_id` ASC),
ADD  CONSTRAINT `movement_ibfk_8`
FOREIGN KEY (`pallet_id`)
REFERENCES `item` (`id`)
ON DELETE CASCADE ON UPDATE NO ACTION;

ALTER TABLE `stock_movement`
ADD COLUMN `is_not_full` TINYINT(1) NOT NULL DEFAULT 0;