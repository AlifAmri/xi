SET FOREIGN_KEY_CHECKS = 0;
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