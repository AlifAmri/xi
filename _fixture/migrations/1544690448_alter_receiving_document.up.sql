SET FOREIGN_KEY_CHECKS = 0;
ALTER TABLE `receiving_document`
DROP FOREIGN KEY `fk_receiving_plan_4`,
ADD FOREIGN KEY (`unit_id`) REFERENCES `stock_unit` (`id`) ON DELETE SET NULL ON UPDATE NO ACTION
