SET FOREIGN_KEY_CHECKS = 0;
INSERT INTO `privilege` (`id`,`name`,`action`,`is_active`,`note`) VALUES
('29', 'report', 'ban', '1', 'Dapat melihat report ban'),
('30', 'report', 'sticker', '1', 'Dapat melihat report sticker'),
('31', 'report', 'pallet', '1', 'Dapat melihat report pallet'),
('32', 'report', 'spacer', '1', 'Dapat melihat report spacer');

UPDATE `privilege` SET action='storage', note='Dapat melihat report storage' WHERE id=21 ;