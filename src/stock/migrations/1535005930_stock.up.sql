SET FOREIGN_KEY_CHECKS = 0;
CREATE TABLE `stock_storage` (
  `id` BIGINT(11) UNSIGNED NOT NULL AUTO_INCREMENT,
  `location_id` BIGINT(11) UNSIGNED NULL,
  `container_id` BIGINT(11) UNSIGNED NULL,
  `code` VARCHAR(45) NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  INDEX `fk_stock_storage_1_idx` (`location_id` ASC),
  INDEX `fk_stock_storage_2_idx` (`container_id` ASC),
  CONSTRAINT `fk_stock_storage_1_idx`
    FOREIGN KEY (`location_id`)
    REFERENCES `warehouse_location` (`id`)
    ON DELETE SET NULL
    ON UPDATE NO ACTION,
  CONSTRAINT `fk_stock_storage_2_idx`
    FOREIGN KEY (`container_id`)
    REFERENCES `item` (`id`)
    ON DELETE SET NULL
    ON UPDATE NO ACTION);

CREATE TABLE `stock_unit` (
  `id` BIGINT(11) UNSIGNED NOT NULL AUTO_INCREMENT,
  `item_id` BIGINT(11) UNSIGNED NOT NULL,
  `batch_id` BIGINT(11) UNSIGNED NULL,
  `storage_id` BIGINT(11) UNSIGNED NULL,
  `ref_id` BIGINT(11) UNSIGNED NULL,
  `code` VARCHAR(45) NOT NULL,
  `stock` BIGINT(11) NOT NULL DEFAULT 0,
  `is_defect` TINYINT(1) NOT NULL DEFAULT 0,
  `status` ENUM('draft', 'stored', 'moving', 'prepared', 'void', 'out') NOT NULL DEFAULT 'draft',
  `barcode_image` TINYTEXT NULL,
  `created_by` bigint(11) unsigned NULL,
  `received_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `stock_unit_identifier` (`code` ASC),
  INDEX `fk_stock_unit_1_idx` (`item_id` ASC),
  INDEX `fk_stock_unit_2_idx` (`batch_id` ASC),
  INDEX `fk_stock_unit_3_idx` (`storage_id` ASC),
  INDEX `fk_stock_unit_4_idx` (`created_by` ASC),
  CONSTRAINT `fk_stock_unit_1`
    FOREIGN KEY (`item_id`)
    REFERENCES `item` (`id`)
    ON DELETE CASCADE
    ON UPDATE NO ACTION,
  CONSTRAINT `fk_stock_unit_2`
    FOREIGN KEY (`batch_id`)
    REFERENCES `item_batch` (`id`)
    ON DELETE CASCADE
    ON UPDATE NO ACTION,
  CONSTRAINT `fk_stock_unit_3`
    FOREIGN KEY (`storage_id`)
    REFERENCES `stock_storage` (`id`)
    ON DELETE SET NULL
    ON UPDATE NO ACTION,
  CONSTRAINT `fk_stock_unit_4`
    FOREIGN KEY (`created_by`)
    REFERENCES `user` (`id`)
    ON DELETE SET NULL
    ON UPDATE NO ACTION);

CREATE TABLE `stock_movement` (
  `id` bigint(11) unsigned NOT NULL AUTO_INCREMENT,
  `unit_id` bigint(11) unsigned NOT NULL,
  `code` VARCHAR(45) NOT NULL,
  `type` enum('routine','picking','putaway') NOT NULL DEFAULT 'routine',
  `ref_id` BIGINT(11) UNSIGNED DEFAULT NULL,
  `ref_code` VARCHAR(45) DEFAULT NULL,
  `status` enum('new','start','finish') NOT NULL DEFAULT 'new',
  `quantity` decimal(12,2) unsigned NOT NULL DEFAULT '0.00',
  `is_partial` tinyint(1) unsigned NOT NULL DEFAULT '0',
  `is_merger` tinyint(1) unsigned NOT NULL DEFAULT '0',
  `origin_id` bigint(11) unsigned DEFAULT NULL,
  `destination_id` bigint(11) unsigned DEFAULT NULL,
  `new_unit` bigint(11) unsigned DEFAULT NULL,
  `merge_unit` bigint(11) unsigned DEFAULT NULL,
  `note` tinytext,
  `created_by` bigint(11) unsigned NOT NULL,
  `moved_by` bigint(11) unsigned DEFAULT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `started_at` timestamp NULL DEFAULT NULL,
  `finished_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `fk_movement_1_idx` (`created_by`),
  KEY `unit_id` (`unit_id`),
  KEY `origin_id` (`origin_id`),
  KEY `destination_id` (`destination_id`),
  KEY `new_unit` (`new_unit`),
  KEY `merge_unit` (`merge_unit`),
  KEY `moved_by` (`moved_by`),
  CONSTRAINT `fk_movement_1` FOREIGN KEY (`created_by`) REFERENCES `user` (`id`) ON DELETE NO ACTION ON UPDATE NO ACTION,
  CONSTRAINT `movement_ibfk_1` FOREIGN KEY (`unit_id`) REFERENCES `stock_unit` (`id`),
  CONSTRAINT `movement_ibfk_2` FOREIGN KEY (`origin_id`) REFERENCES `warehouse_location` (`id`) ON DELETE CASCADE,
  CONSTRAINT `movement_ibfk_3` FOREIGN KEY (`destination_id`) REFERENCES `warehouse_location` (`id`) ON DELETE CASCADE,
  CONSTRAINT `movement_ibfk_4` FOREIGN KEY (`new_unit`) REFERENCES `stock_unit` (`id`) ON DELETE CASCADE,
  CONSTRAINT `movement_ibfk_5` FOREIGN KEY (`merge_unit`) REFERENCES `stock_unit` (`id`) ON DELETE CASCADE,
  CONSTRAINT `movement_ibfk_6` FOREIGN KEY (`created_by`) REFERENCES `user` (`id`) ON DELETE CASCADE,
  CONSTRAINT `movement_ibfk_7` FOREIGN KEY (`moved_by`) REFERENCES `user` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

CREATE TABLE `stock_opname` (
  `id` BIGINT(11) UNSIGNED NOT NULL AUTO_INCREMENT,
  `location_id` BIGINT(11) UNSIGNED NULL,
  `code` VARCHAR(45) NOT NULL,
  `type` ENUM('opname', 'adjustment') NOT NULL DEFAULT 'opname',
  `status` ENUM('active', 'finish', 'cancelled') NOT NULL DEFAULT 'active',
  `note` TINYTEXT NULL,
  `created_by` BIGINT(11) UNSIGNED NULL,
  `approved_by` BIGINT(11) UNSIGNED NULL,
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `approved_at` TIMESTAMP NULL,
  PRIMARY KEY (`id`),
  INDEX `fk_stockopname_1_idx` (`location_id` ASC),
  INDEX `fk_stockopname_2_idx` (`created_by` ASC),
  INDEX `fk_stockopname_3_idx` (`approved_by` ASC),
  CONSTRAINT `fk_stockopname_1`
    FOREIGN KEY (`location_id`)
    REFERENCES `warehouse_location` (`id`)
    ON DELETE SET NULL
    ON UPDATE NO ACTION,
  CONSTRAINT `fk_stockopname_2`
    FOREIGN KEY (`created_by`)
    REFERENCES `user` (`id`)
    ON DELETE SET NULL
    ON UPDATE NO ACTION,
  CONSTRAINT `fk_stockopname_3`
    FOREIGN KEY (`approved_by`)
    REFERENCES `user` (`id`)
    ON DELETE SET NULL
    ON UPDATE NO ACTION);

CREATE TABLE `stock_opname_item` (
  `id` BIGINT(11) UNSIGNED NOT NULL AUTO_INCREMENT,
  `stock_opname_id` BIGINT(11) UNSIGNED NOT NULL,
  `item_id` BIGINT(11) UNSIGNED NULL DEFAULT NULL,
  `unit_id` BIGINT(11) UNSIGNED NULL DEFAULT NULL,
  `container_id` BIGINT(11) UNSIGNED NULL DEFAULT NULL,
  `container_num` TINYINT(2) NOT NULL DEFAULT 0,
  `unit_quantity` DECIMAL(12,2) NOT NULL DEFAULT 0,
  `actual_quantity` DECIMAL(12,2) NOT NULL DEFAULT 0,
  `is_new_unit` TINYINT(1) NOT NULL DEFAULT 0,
  `is_defect` TINYINT(1) NOT NULL DEFAULT 0,
  `is_void` TINYINT(1) NOT NULL DEFAULT 0,
  `note` TINYTEXT NULL,
  PRIMARY KEY (`id`),
  INDEX `fk_stock_opname_item_1_idx` (`stock_opname_id` ASC),
  INDEX `fk_stock_opname_item_2_idx` (`item_id` ASC),
  INDEX `fk_stock_opname_item_3_idx` (`unit_id` ASC),
  INDEX `fk_stock_opname_item_4_idx` (`container_id` ASC),
  CONSTRAINT `fk_stock_opname_item_1`
    FOREIGN KEY (`stock_opname_id`)
    REFERENCES `stock_opname` (`id`)
    ON DELETE CASCADE
    ON UPDATE NO ACTION,
  CONSTRAINT `fk_stock_opname_item_2`
    FOREIGN KEY (`item_id`)
    REFERENCES `item` (`id`)
    ON DELETE CASCADE
    ON UPDATE NO ACTION,
  CONSTRAINT `fk_stock_opname_item_4`
    FOREIGN KEY (`container_id`)
    REFERENCES `item` (`id`)
    ON DELETE CASCADE
    ON UPDATE NO ACTION,
  CONSTRAINT `fk_stock_opname_item_3`
    FOREIGN KEY (`unit_id`)
    REFERENCES `stock_unit` (`id`)
    ON DELETE CASCADE
    ON UPDATE NO ACTION);

CREATE TABLE `stock_log` (
  `id` BIGINT(11) UNSIGNED NOT NULL AUTO_INCREMENT,
  `stock_unit_id` BIGINT(11) UNSIGNED NULL,
  `item_id` BIGINT(11) UNSIGNED NULL,
  `batch_id` BIGINT(11) UNSIGNED NULL,
  `ref_type` VARCHAR(45) NULL,
  `ref_id` BIGINT(11) UNSIGNED NULL,
  `ref_code` VARCHAR(45) NOT NULL,
  `quantity` DECIMAL(20,2) NULL,
  `recorded_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  INDEX `fk_stock_log_1_idx` (`stock_unit_id` ASC),
  INDEX `fk_stock_log_2_idx` (`item_id` ASC),
  INDEX `fk_stock_log_3_idx` (`batch_id` ASC),
  CONSTRAINT `fk_stock_log_1`
    FOREIGN KEY (`stock_unit_id`)
    REFERENCES `stock_unit` (`id`)
    ON DELETE CASCADE
    ON UPDATE NO ACTION,
  CONSTRAINT `fk_stock_log_2`
    FOREIGN KEY (`item_id`)
    REFERENCES `item` (`id`)
    ON DELETE CASCADE
    ON UPDATE NO ACTION,
  CONSTRAINT `fk_stock_log_3`
    FOREIGN KEY (`batch_id`)
    REFERENCES `item_batch` (`id`)
    ON DELETE SET NULL
    ON UPDATE NO ACTION);
