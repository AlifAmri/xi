SET FOREIGN_KEY_CHECKS = 0;
ALTER TABLE `receiving_unit`
ADD COLUMN  `pallet_id` BIGINT(11) UNSIGNED NOT NULL,
ADD INDEX `fk_receiving_unit_9_idx` (`pallet_id` ASC),
ADD  CONSTRAINT `fk_receiving_unit_9`
FOREIGN KEY (`pallet_id`)
REFERENCES `item` (`id`)
ON DELETE CASCADE ON UPDATE NO ACTION;

ALTER TABLE `receiving_unit`
ADD COLUMN `is_not_full` TINYINT(1) NOT NULL DEFAULT 0;