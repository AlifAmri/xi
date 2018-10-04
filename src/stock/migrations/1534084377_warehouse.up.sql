SET FOREIGN_KEY_CHECKS = 0;

CREATE TABLE `warehouse` (
  `id` BIGINT(11) UNSIGNED NOT NULL AUTO_INCREMENT,
  `code` VARCHAR(45) NOT NULL,
  `name` VARCHAR(45) NOT NULL,
  `note` TINYTEXT NULL,
  PRIMARY KEY (`id`));

CREATE TABLE `warehouse_area` (
  `id` BIGINT(11) UNSIGNED NOT NULL AUTO_INCREMENT,
  `warehouse_id` BIGINT(11) UNSIGNED NOT NULL,
  `code` VARCHAR(45) NOT NULL,
  `name` VARCHAR(45) NOT NULL,
  `is_active` TINYINT(1) NOT NULL DEFAULT 1,
  `note` TINYTEXT NULL,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `code_UNIQUE` (`warehouse_id` ASC, `code` ASC),
  CONSTRAINT `fk_warehouse_area_1`
    FOREIGN KEY (`warehouse_id`)
    REFERENCES `warehouse` (`id`)
    ON DELETE CASCADE
    ON UPDATE NO ACTION);

CREATE TABLE `warehouse_location` (
  `id` BIGINT(11) UNSIGNED NOT NULL AUTO_INCREMENT,
  `warehouse_area_id` BIGINT(11) UNSIGNED NOT NULL,
  `code` VARCHAR(45) NOT NULL,
  `name` VARCHAR(45) NOT NULL,
  `is_active` TINYINT(1) NOT NULL DEFAULT 1,
  `coordinate_x` INT(11) NOT NULL,
  `coordinate_y` INT(11) NOT NULL,
  `coordinate_w` INT(11) NOT NULL,
  `coordinate_h` INT(11) NOT NULL,
  `note` TINYTEXT NULL,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `index2` (`warehouse_area_id` ASC, `code` ASC),
  CONSTRAINT `fk_warehouse_location_1`
    FOREIGN KEY (`warehouse_area_id`)
    REFERENCES `warehouse_area` (`id`)
    ON DELETE CASCADE
    ON UPDATE NO ACTION);