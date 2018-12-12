SET FOREIGN_KEY_CHECKS = 0;
ALTER TABLE `preparation_unit`
ADD COLUMN  `qc_by` BIGINT(11) UNSIGNED NULL,
ADD INDEX `fk_preparation_unit_7_idx` (`qc_by` ASC),
ADD  CONSTRAINT `preparation_unit_7_idx`
FOREIGN KEY (`qc_by`)
REFERENCES `user` (`id`)
ON DELETE CASCADE ON UPDATE NO ACTION;