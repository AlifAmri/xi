SET FOREIGN_KEY_CHECKS = 0;

CREATE TABLE `item_type` (
  `id` BIGINT(11) UNSIGNED NOT NULL AUTO_INCREMENT,
  `name` VARCHAR(45) NOT NULL,
  `is_batch` TINYINT(1) NOT NULL DEFAULT 1,
  `note` TINYTEXT NULL,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `name_UNIQUE` (`name` ASC));

CREATE TABLE `item_group` (
  `id` BIGINT(11) UNSIGNED NOT NULL AUTO_INCREMENT,
  `type_id` BIGINT(11) UNSIGNED NULL DEFAULT NULL,
  `name` VARCHAR(100) NOT NULL,
  `note` TINYTEXT NULL,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `name_UNIQUE` (`type_id` ASC, `name` ASC),
  INDEX `fk_item_group_type_idx` (`type_id` ASC),
  CONSTRAINT `fk_item_group_type_idx` FOREIGN KEY (`type_id`) REFERENCES `item_type` (`id`) ON DELETE SET NULL ON UPDATE NO ACTION);

CREATE TABLE `item_group_attribute` (
  `id` BIGINT(11) UNSIGNED NOT NULL AUTO_INCREMENT,
  `item_group_id` BIGINT(11) UNSIGNED NOT NULL,
  `attribute` VARCHAR(45) NOT NULL,
  `value` TINYTEXT NULL,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `group_attribute` (`item_group_id` ASC, `attribute` ASC),
  CONSTRAINT `fk_item_group_attribute_1`
    FOREIGN KEY (`item_group_id`)
    REFERENCES `item_group` (`id`)
    ON DELETE CASCADE
    ON UPDATE NO ACTION);

CREATE TABLE `item_category` (
  `id` BIGINT(11) UNSIGNED NOT NULL AUTO_INCREMENT,
  `type_id` BIGINT(11) UNSIGNED NULL DEFAULT NULL,
  `parent_id` BIGINT(11) NULL DEFAULT NULL,
  `name` VARCHAR(45) NOT NULL,
  `note` VARCHAR(45) NULL,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `name_UNIQUE` (`type_id` ASC, `name` ASC),
  INDEX `fk_item_category_type_idx` (`type_id` ASC),
  CONSTRAINT `fk_item_category_type_idx` FOREIGN KEY (`type_id`) REFERENCES `item_type` (`id`) ON DELETE SET NULL ON UPDATE NO ACTION);

CREATE TABLE `item_uom` (
  `id` BIGINT(11) UNSIGNED NOT NULL AUTO_INCREMENT,
  `type_id` BIGINT(11) UNSIGNED NULL DEFAULT NULL,
  `name` VARCHAR(45) NOT NULL,
  `note` TINYTEXT NULL,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `name_UNIQUE` (`type_id` ASC, `name` ASC),
  INDEX `fk_item_uom_type_idx` (`type_id` ASC),
  CONSTRAINT `fk_item_uom_type_idx` FOREIGN KEY (`type_id`) REFERENCES `item_type` (`id`) ON DELETE SET NULL ON UPDATE NO ACTION);

CREATE TABLE `item` (
  `id` BIGINT(11) UNSIGNED NOT NULL AUTO_INCREMENT,
  `group_id` BIGINT(11) UNSIGNED NULL DEFAULT NULL,
  `type_id` BIGINT(11) UNSIGNED NULL DEFAULT NULL,
  `category_id` BIGINT(11) UNSIGNED NULL DEFAULT NULL,
  `preferred_area` BIGINT(11) UNSIGNED NULL DEFAULT NULL,
  `code` VARCHAR(45) NOT NULL,
  `name` VARCHAR(145) NOT NULL,
  `stock` BIGINT(11) UNSIGNED NOT NULL DEFAULT 0,
  `is_active` TINYINT(1) NOT NULL DEFAULT 1,
  `image` TINYTEXT NULL,
  `barcode_number` VARCHAR(50) NULL,
  `barcode_image` TINYTEXT NULL,
  `note` TINYTEXT NULL,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `item_identifier` (`type_id` ASC, `code` ASC),
  INDEX `fk_item_1_idx` (`group_id` ASC),
  INDEX `fk_item_4_idx` (`preferred_area` ASC),
  INDEX `fk_item_3_idx` (`category_id` ASC),
  INDEX `fk_item_2_idx` (`type_id` ASC),
  CONSTRAINT `fk_item_1`
    FOREIGN KEY (`group_id`)
    REFERENCES `item_group` (`id`)
    ON DELETE SET NULL
    ON UPDATE NO ACTION,
  CONSTRAINT `fk_item_2`
    FOREIGN KEY (`type_id`)
    REFERENCES `item_type` (`id`)
    ON DELETE SET NULL
    ON UPDATE NO ACTION,
  CONSTRAINT `fk_item_3`
    FOREIGN KEY (`category_id`)
    REFERENCES `item_category` (`id`)
    ON DELETE SET NULL
    ON UPDATE NO ACTION,
  CONSTRAINT `fk_item_4`
    FOREIGN KEY (`preferred_area`)
    REFERENCES `warehouse_area` (`id`)
    ON DELETE SET NULL
    ON UPDATE NO ACTION);

CREATE TABLE `item_attribute` (
  `id` BIGINT(11) UNSIGNED NOT NULL AUTO_INCREMENT,
  `item_id` BIGINT(11) UNSIGNED NOT NULL,
  `attribute` VARCHAR(45) NULL,
  `value` TINYTEXT NULL,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `item_attribute` (`item_id` ASC, `attribute` ASC),
  CONSTRAINT `fk_item_attribute_1`
    FOREIGN KEY (`item_id`)
    REFERENCES `item` (`id`)
    ON DELETE CASCADE
    ON UPDATE NO ACTION);

CREATE TABLE `item_batch` (
  `id` BIGINT(11) UNSIGNED NOT NULL AUTO_INCREMENT,
  `item_id` BIGINT(11) UNSIGNED NOT NULL,
  `code` VARCHAR(45) NOT NULL,
  `name` VARCHAR(45) NULL DEFAULT NULL,
  `stock` BIGINT(11) UNSIGNED NOT NULL DEFAULT 0,
  `produced_at` TIMESTAMP NULL,
  `expired_at` TIMESTAMP NULL,
  `entry_at` TIMESTAMP NULL,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `batch_identifier` (`code` ASC),
  INDEX `fk_item_batch_1_idx` (`item_id` ASC),
  CONSTRAINT `fk_item_batch_1`
    FOREIGN KEY (`item_id`)
    REFERENCES `item` (`id`)
    ON DELETE CASCADE
    ON UPDATE NO ACTION);
