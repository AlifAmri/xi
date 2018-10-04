SET FOREIGN_KEY_CHECKS = 0;

CREATE TABLE `app_config` (
  `id` BIGINT(11) UNSIGNED NOT NULL AUTO_INCREMENT,
  `attribute` VARCHAR(45) NOT NULL,
  `value` TINYTEXT NULL,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `config_identifier` (`attribute` ASC));

CREATE TABLE `privilege` (
`id` BIGINT(11) UNSIGNED NOT NULL AUTO_INCREMENT,
`name` VARCHAR(45) NOT NULL,
`action` VARCHAR(45) NOT NULL,
`is_active` TINYINT(1) UNSIGNED NOT NULL DEFAULT 1,
`note` TINYTEXT NULL,
PRIMARY KEY (`id`),
UNIQUE INDEX `privilege_module` (`name` ASC, `action` ASC));

CREATE TABLE `privilege_user` (
  `id` BIGINT(11) UNSIGNED NOT NULL AUTO_INCREMENT,
  `privilege_id` BIGINT(11) UNSIGNED NOT NULL,
  `user_id` BIGINT(11) UNSIGNED NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `user_permission` (`privilege_id` ASC, `user_id` ASC),
  CONSTRAINT `fk_privilege_user_1`
  FOREIGN KEY (`privilege_id`)
  REFERENCES `privilege` (`id`) ON DELETE CASCADE ON UPDATE NO ACTION);

CREATE TABLE `privilege_usergroup` (
  `id` BIGINT(11) UNSIGNED NOT NULL AUTO_INCREMENT,
  `privilege_id` BIGINT(11) UNSIGNED NOT NULL,
  `usergroup_id` BIGINT(11) UNSIGNED NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `usergroup_permision` (`privilege_id` ASC, `usergroup_id` ASC),
  CONSTRAINT `fk_privilege_usergroup_1`
  FOREIGN KEY (`privilege_id`)
  REFERENCES `privilege` (`id`) ON DELETE CASCADE ON UPDATE NO ACTION);

CREATE TABLE `usergroup` (
  `id` BIGINT(11) UNSIGNED NOT NULL AUTO_INCREMENT,
  `name` VARCHAR(45) NOT NULL,
  `note` TINYTEXT NULL,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `usergroup_unique` (`name` ASC));

CREATE TABLE `user` (
  `id` BIGINT(11) UNSIGNED NOT NULL AUTO_INCREMENT,
  `usergroup_id` BIGINT(11) UNSIGNED NULL DEFAULT  NULL,
  `username` VARCHAR(50) NOT NULL,
  `password` VARCHAR(150) NOT NULL,
  `name` VARCHAR(80) NOT NULL,
  `is_active` TINYINT(1) NOT NULL DEFAULT 1,
  `is_superuser` TINYINT(1) NOT NULL DEFAULT 0,
  `last_login_at` TIMESTAMP NULL,
  `registered_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `user_identification` (`username` ASC),
  INDEX `fk_user_1_idx` (`usergroup_id` ASC),
  CONSTRAINT `fk_user_1`
    FOREIGN KEY (`usergroup_id`)
    REFERENCES `usergroup` (`id`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION);

CREATE TABLE `warehouse` (
  `id` BIGINT(11) UNSIGNED NOT NULL AUTO_INCREMENT,
  `code` VARCHAR(45) NOT NULL,
  `name` VARCHAR(45) NOT NULL,
  `note` TINYTEXT NULL,
  PRIMARY KEY (`id`));

CREATE TABLE `warehouse_area` (
  `id` BIGINT(11) UNSIGNED NOT NULL AUTO_INCREMENT,
  `warehouse_id` BIGINT(11) UNSIGNED NOT NULL,
  `type` ENUM('storage', 'preparation', 'receiving', 'other') NOT NULL DEFAULT 'storage',
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
  `storage_capacity` TINYINT(2) NOT NULL DEFAULT 0,
  `storage_used` TINYINT(2) NOT NULL DEFAULT 0,
  `note` TINYTEXT NULL,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `index2` (`warehouse_area_id` ASC, `code` ASC),
  CONSTRAINT `fk_warehouse_location_1`
    FOREIGN KEY (`warehouse_area_id`)
    REFERENCES `warehouse_area` (`id`)
    ON DELETE CASCADE
    ON UPDATE NO ACTION);

CREATE TABLE `partnership_type` (
  `id` BIGINT(11) UNSIGNED NOT NULL AUTO_INCREMENT,
  `name` VARCHAR(45) NOT NULL,
  `note` TINYTEXT NULL,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `name_UNIQUE` (`name` ASC));

CREATE TABLE `partnership` (
  `id` BIGINT(11) UNSIGNED NOT NULL AUTO_INCREMENT,
  `type_id` BIGINT(11) UNSIGNED NOT NULL,
  `company_name` VARCHAR(100) NULL,
  `company_address` TINYTEXT NULL,
  `company_phone` VARCHAR(45) NULL,
  `company_email` VARCHAR(100) NULL,
  `contact_person` VARCHAR(100) NULL,
  `is_active` TINYINT(1) NOT NULL DEFAULT 1,
  `note` TINYTEXT NULL,
  PRIMARY KEY (`id`),
  INDEX `fk_partnership_1_idx` (`type_id` ASC),
  CONSTRAINT `fk_partnership_1`
    FOREIGN KEY (`type_id`)
    REFERENCES `partnership_type` (`id`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION);

CREATE TABLE `item_type` (
  `id` BIGINT(11) UNSIGNED NOT NULL AUTO_INCREMENT,
  `name` VARCHAR(45) NOT NULL,
  `is_batch` TINYINT(1) NOT NULL DEFAULT 1,
  `is_container` TINYINT(1) NOT NULL DEFAULT 1,
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
  `stock` DECIMAL(12,2) NOT NULL DEFAULT 0,
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
  `stock` DECIMAL(12,2) NOT NULL DEFAULT 0,
  `produced_at` TIMESTAMP NULL,
  `expired_at` TIMESTAMP NULL,
  `entry_at` TIMESTAMP NULL,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `batch_identifier` (`code`, `item_id` ASC),
  INDEX `fk_item_batch_1_idx` (`item_id` ASC),
  CONSTRAINT `fk_item_batch_1`
    FOREIGN KEY (`item_id`)
    REFERENCES `item` (`id`)
    ON DELETE CASCADE
    ON UPDATE NO ACTION);

CREATE TABLE `storage_group` (
  `id` BIGINT(11) UNSIGNED NOT NULL AUTO_INCREMENT,
  `name` VARCHAR(45) NOT NULL,
  `type` ENUM('default', 'item_group', 'item_category', 'item_batch', 'item_code', 'ncp') NOT NULL DEFAULT 'default',
  `type_value` VARCHAR(45) NULL DEFAULT NULL,
  `is_active` TINYINT(1) NOT NULL DEFAULT 0,
  `is_primary` TINYINT(1) NOT NULL DEFAULT 0,
  `note` TINYTEXT NULL,
  PRIMARY KEY (`id`));

CREATE TABLE `storage_group_area` (
  `id` BIGINT(11) UNSIGNED NOT NULL AUTO_INCREMENT,
  `storage_group_id` BIGINT(11) UNSIGNED NOT NULL,
  `warehouse_area_id` BIGINT(11) UNSIGNED NOT NULL,
  PRIMARY KEY (`id`),
  INDEX `fk_storage_group_area_1_idx` (`storage_group_id` ASC),
  INDEX `fk_storage_group_area_2_idx` (`warehouse_area_id` ASC),
  CONSTRAINT `fk_storage_group_area_1`
    FOREIGN KEY (`storage_group_id`)
    REFERENCES `storage_group` (`id`)
    ON DELETE CASCADE
    ON UPDATE NO ACTION,
  CONSTRAINT `fk_storage_group_area_2`
    FOREIGN KEY (`warehouse_area_id`)
    REFERENCES `warehouse_area` (`id`)
    ON DELETE CASCADE
    ON UPDATE NO ACTION);

CREATE TABLE `incoming_vehicle` (
  `id` BIGINT(11) UNSIGNED NOT NULL AUTO_INCREMENT,
  `document_id` BIGINT(11) UNSIGNED NULL DEFAULT NULL,
  `purpose` ENUM('receiving', 'dispatching', 'other') NOT NULL DEFAULT 'receiving',
  `status` ENUM('in_progress', 'finished', 'out') NOT NULL DEFAULT 'in_progress',
  `vehicle_type` VARCHAR(45) NOT NULL,
  `vehicle_number` VARCHAR(45) NOT NULL,
  `driver` VARCHAR(145) NOT NULL,
  `picture` TINYTEXT NULL,
  `cargo_type` ENUM('curah', 'pallet', 'mix') NOT NULL DEFAULT 'curah',
  `subcon_id` BIGINT(11) UNSIGNED NULL,
  `subcon_typed` VARCHAR(45) NULL,
  `destination` ENUM('local', 'export', 'import') NOT NULL DEFAULT 'local',
  `container_number` VARCHAR(45) NULL,
  `seal_number` VARCHAR(45) NULL,
  `notes` TINYTEXT NULL,
  `in_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `out_at` TIMESTAMP NULL,
  PRIMARY KEY (`id`),
  INDEX `fk_incoming vehicles_1_idx` (`subcon_id` ASC),
  CONSTRAINT `fk_incoming vehicles_1`
    FOREIGN KEY (`subcon_id`)
    REFERENCES `partnership` (`id`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION);

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

CREATE TABLE `receipt_plan` (
  `id` BIGINT(11) UNSIGNED NOT NULL AUTO_INCREMENT,
  `partner_id` BIGINT(11) UNSIGNED NULL,
  `status` ENUM('draft', 'pending', 'active', 'finish') NOT NULL DEFAULT 'draft',
  `document_code` VARCHAR(45) NOT NULL,
  `total_quantity` DECIMAL(12,2) NOT NULL DEFAULT 0,
  `note` TINYTEXT NULL,
  `approved_by` BIGINT(11) UNSIGNED NULL,
  `created_by` BIGINT(11) UNSIGNED NOT NULL,
  `created_at` TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP,
  `received_at` TIMESTAMP NULL,
  PRIMARY KEY (`id`),
  INDEX `fk_receipt_plan_1_idx` (`partner_id` ASC),
  INDEX `fk_receipt_plan_2_idx` (`approved_by` ASC),
  INDEX `fk_receipt_plan_3_idx` (`created_by` ASC),
  CONSTRAINT `fk_receipt_plan_1`
    FOREIGN KEY (`partner_id`)
    REFERENCES `partnership` (`id`)
    ON DELETE SET NULL
    ON UPDATE NO ACTION,
  CONSTRAINT `fk_receipt_plan_2`
    FOREIGN KEY (`approved_by`)
    REFERENCES `user` (`id`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION,
  CONSTRAINT `fk_receipt_plan_3`
    FOREIGN KEY (`created_by`)
    REFERENCES `user` (`id`)
    ON DELETE CASCADE
    ON UPDATE NO ACTION);

CREATE TABLE `receipt_plan_item` (
  `id` bigint(11) unsigned NOT NULL AUTO_INCREMENT,
  `plan_id` bigint(11) unsigned NOT NULL,
  `unit_code` varchar(100) DEFAULT NULL,
  `item_code` varchar(100) NOT NULL,
  `batch_code` varchar(100) DEFAULT NULL,
  `quantity` decimal(12,2) NOT NULL DEFAULT '0.00',
  PRIMARY KEY (`id`),
  KEY `fk_plan_item_1_idx` (`plan_id`),
  CONSTRAINT `fk_plan_item_1` FOREIGN KEY (`plan_id`) REFERENCES `receipt_plan` (`id`) ON DELETE CASCADE ON UPDATE NO ACTION);

CREATE TABLE `receiving` (
  `id` BIGINT(11) UNSIGNED NOT NULL AUTO_INCREMENT,
  `vehicle_id` BIGINT(11) UNSIGNED NOT NULL,
  `plan_id` BIGINT(11) UNSIGNED NULL,
  `partner_id` BIGINT(11) UNSIGNED NULL,
  `supervisor_id` BIGINT(11) UNSIGNED NULL,
  `code` VARCHAR(45) NULL,
  `status` ENUM('active', 'finish') NOT NULL DEFAULT 'active',
  `document_code` VARCHAR(45) NULL,
  `document_file` TINYTEXT NULL,
  `total_quantity_plan` DECIMAL(12,2) NULL,
  `total_quantity_actual` DECIMAL(12,2) NULL,
  `note` TINYTEXT NULL,
  `approved_by` BIGINT(11) UNSIGNED NULL,
  `created_by` BIGINT(11) UNSIGNED NULL,
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `received_at` TIMESTAMP NULL,
  `started_at` TIMESTAMP NULL,
  `finished_at` TIMESTAMP NULL,
  PRIMARY KEY (`id`),
  INDEX `fk_receiving_1_idx` (`vehicle_id` ASC),
  INDEX `fk_receiving_2_idx` (`plan_id` ASC),
  INDEX `fk_receiving_3_idx` (`partner_id` ASC),
  INDEX `fk_receiving_4_idx` (`supervisor_id` ASC),
  INDEX `fk_receiving_5_idx` (`approved_by` ASC),
  INDEX `fk_receiving_6_idx` (`created_by` ASC),
  CONSTRAINT `fk_receiving_1`
    FOREIGN KEY (`vehicle_id`)
    REFERENCES `incoming_vehicle` (`id`)
    ON DELETE CASCADE
    ON UPDATE NO ACTION,
  CONSTRAINT `fk_receiving_2`
    FOREIGN KEY (`plan_id`)
    REFERENCES `receipt_plan` (`id`)
    ON DELETE SET NULL
    ON UPDATE NO ACTION,
  CONSTRAINT `fk_receiving_3`
    FOREIGN KEY (`partner_id`)
    REFERENCES `partnership` (`id`)
    ON DELETE SET NULL
    ON UPDATE NO ACTION,
  CONSTRAINT `fk_receiving_4`
    FOREIGN KEY (`supervisor_id`)
    REFERENCES `user` (`id`)
    ON DELETE SET NULL
    ON UPDATE NO ACTION,
  CONSTRAINT `fk_receiving_5`
    FOREIGN KEY (`approved_by`)
    REFERENCES `user` (`id`)
    ON DELETE SET NULL
    ON UPDATE NO ACTION,
  CONSTRAINT `fk_receiving_6`
    FOREIGN KEY (`created_by`)
    REFERENCES `user` (`id`)
    ON DELETE SET NULL
    ON UPDATE NO ACTION);


CREATE TABLE `receiving_document` (
  `id` BIGINT(11) UNSIGNED NOT NULL AUTO_INCREMENT,
  `receiving_id` bigint(11) UNSIGNED NOT NULL,
  `item_id` BIGINT(11) UNSIGNED NOT NULL,
  `batch_id` BIGINT(11) UNSIGNED DEFAULT NULL,
  `unit_id` BIGINT(11) UNSIGNED DEFAULT NULL,
  `is_new` TINYINT(1) NOT NULL DEFAULT '0',
  `quantity` decimal(12,2) NOT NULL DEFAULT '0.00',
  PRIMARY KEY (`id`),
  KEY `fk_receiving_plan_1_idx` (`receiving_id`),
  KEY `fk_receiving_plan_2_idx` (`item_id`),
  KEY `fk_receiving_plan_3_idx` (`batch_id`),
  KEY `fk_receiving_plan_4_idx` (`unit_id`),
  CONSTRAINT `fk_receiving_plan_1` FOREIGN KEY (`receiving_id`) REFERENCES `receiving` (`id`) ON DELETE CASCADE ON UPDATE NO ACTION,
  CONSTRAINT `fk_receiving_plan_2` FOREIGN KEY (`item_id`) REFERENCES `item` (`id`) ON DELETE NO ACTION ON UPDATE NO ACTION,
  CONSTRAINT `fk_receiving_plan_3` FOREIGN KEY (`batch_id`) REFERENCES `item_batch` (`id`) ON DELETE NO ACTION ON UPDATE NO ACTION,
  CONSTRAINT `fk_receiving_plan_4` FOREIGN KEY (`unit_id`) REFERENCES `stock_unit` (`id`) ON DELETE NO ACTION ON UPDATE NO ACTION);

CREATE TABLE `receiving_unit` (
  `id` BIGINT(11) UNSIGNED NOT NULL AUTO_INCREMENT,
  `receiving_id` BIGINT(11) UNSIGNED NOT NULL,
  `unit_id` BIGINT(11) UNSIGNED NULL,
  `location_received` BIGINT(11) UNSIGNED NULL,
  `location_suggested` BIGINT(11) UNSIGNED NULL,
  `location_moved` BIGINT(11) UNSIGNED NULL,
  `unit_code` varchar(100) DEFAULT NULL,
  `item_code` varchar(100) DEFAULT NULL,
  `batch_code` varchar(100) DEFAULT NULL,
  `quantity` DECIMAL(12,2) NOT NULL,
  `is_ncp` TINYINT(1) NOT NULL DEFAULT 0,
  `is_active` TINYINT(1) NOT NULL DEFAULT 0,
  `checked_by` BIGINT(11) UNSIGNED NULL,
  `created_by` BIGINT(11) UNSIGNED NULL,
  `approved_by` BIGINT(11) UNSIGNED NULL,
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  INDEX `fk_receiving_unit_1_idx` (`receiving_id` ASC),
  INDEX `fk_receiving_unit_2_idx` (`unit_id` ASC),
  INDEX `fk_receiving_unit_3_idx` (`location_suggested` ASC),
  INDEX `fk_receiving_unit_4_idx` (`location_moved` ASC),
  INDEX `fk_receiving_unit_5_idx` (`checked_by` ASC),
  INDEX `fk_receiving_unit_6_idx` (`created_by` ASC),
  INDEX `fk_receiving_unit_7_idx` (`approved_by` ASC),
  INDEX `fk_receiving_unit_8_idx` (`location_received` ASC),
  CONSTRAINT `fk_receiving_unit_1`
    FOREIGN KEY (`receiving_id`)
    REFERENCES `receiving` (`id`)
    ON DELETE CASCADE
    ON UPDATE NO ACTION,
  CONSTRAINT `fk_receiving_unit_2`
    FOREIGN KEY (`unit_id`)
    REFERENCES `stock_unit` (`id`)
    ON DELETE CASCADE
    ON UPDATE NO ACTION,
  CONSTRAINT `fk_receiving_unit_3`
    FOREIGN KEY (`location_suggested`)
    REFERENCES `warehouse_location` (`id`)
    ON DELETE SET NULL
    ON UPDATE NO ACTION,
  CONSTRAINT `fk_receiving_unit_4`
    FOREIGN KEY (`location_moved`)
    REFERENCES `warehouse_location` (`id`)
    ON DELETE SET NULL
    ON UPDATE NO ACTION,
  CONSTRAINT `fk_receiving_unit_5`
    FOREIGN KEY (`checked_by`)
    REFERENCES `user` (`id`)
    ON DELETE CASCADE
    ON UPDATE NO ACTION,
  CONSTRAINT `fk_receiving_unit_6`
    FOREIGN KEY (`created_by`)
    REFERENCES `user` (`id`)
    ON DELETE CASCADE
    ON UPDATE NO ACTION,
  CONSTRAINT `fk_receiving_unit_7`
    FOREIGN KEY (`approved_by`)
    REFERENCES `user` (`id`)
    ON DELETE CASCADE
    ON UPDATE NO ACTION,
  CONSTRAINT `fk_receiving_unit_8`
    FOREIGN KEY (`location_received`)
    REFERENCES `warehouse_location` (`id`)
    ON DELETE SET NULL
    ON UPDATE NO ACTION);

CREATE TABLE `receiving_actual` (
  `id` BIGINT(11) UNSIGNED NOT NULL AUTO_INCREMENT,
  `receiving_id` BIGINT(11) UNSIGNED NOT NULL,
  `item_id` BIGINT(11) UNSIGNED NOT NULL,
  `batch_id` BIGINT(11) UNSIGNED NULL,
  `unit_id` BIGINT(11) UNSIGNED NULL,
  `quantity_planned` DECIMAL(12,2) NOT NULL DEFAULT 0,
  `quantity_received` DECIMAL(12,2) NOT NULL DEFAULT 0,
  `quantity_defect` DECIMAL(12,2) NOT NULL DEFAULT 0,
  `note` TINYTEXT NULL,
  PRIMARY KEY (`id`),
  INDEX `fk_receiving_actual_1_idx` (`receiving_id` ASC),
  INDEX `fk_receiving_actual_2_idx` (`item_id` ASC),
  INDEX `fk_receiving_actual_3_idx` (`batch_id` ASC),
  INDEX `fk_receiving_actual_4_idx` (`unit_id` ASC),
  CONSTRAINT `fk_receiving_actual_1`
    FOREIGN KEY (`receiving_id`)
    REFERENCES `receiving` (`id`)
    ON DELETE CASCADE
    ON UPDATE NO ACTION,
  CONSTRAINT `fk_receiving_actual_2`
    FOREIGN KEY (`item_id`)
    REFERENCES `item` (`id`)
    ON DELETE CASCADE
    ON UPDATE NO ACTION,
  CONSTRAINT `fk_receiving_actual_3`
    FOREIGN KEY (`batch_id`)
    REFERENCES `item_batch` (`id`)
    ON DELETE SET NULL
    ON UPDATE NO ACTION,
  CONSTRAINT `fk_receiving_actual_4`
    FOREIGN KEY (`unit_id`)
    REFERENCES `stock_unit` (`id`)
    ON DELETE SET NULL
    ON UPDATE NO ACTION);