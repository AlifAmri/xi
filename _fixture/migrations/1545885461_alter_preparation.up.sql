SET FOREIGN_KEY_CHECKS = 0;
ALTER TABLE `preparation`
ADD COLUMN  `checkout_by` BIGINT(11) UNSIGNED NULL,
ADD INDEX `fk_preparation_8_idx` (`checkout_by` ASC),
ADD  CONSTRAINT `fk_preparation_8`
FOREIGN KEY (`checkout_by`)
REFERENCES `user` (`id`)
ON DELETE SET NULL ON UPDATE NO ACTION;