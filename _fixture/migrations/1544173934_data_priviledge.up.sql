SET FOREIGN_KEY_CHECKS = 0;
INSERT INTO `privilege` (`id`,`name`,`action`,`is_active`,`note`) VALUES
('33', 'receiving', 'kendaraan', '1', 'Dapat melakukan penyelesaian kendaraan receiving'),
('34', 'receiving', 'surat_jalan', '1', 'Dapat membuat surat jalan'),
('35', 'receiving', 'finish', '1', 'Dapat menyelesaikan receiving');

UPDATE `privilege` SET note='Dapat melihat data receiving dan surat jalan' WHERE id=3 ;
UPDATE `privilege` SET note='Dapat melakukan proses receiving put away' WHERE id=5 ;
UPDATE `privilege` SET name='preparation_surat_jalan' WHERE id=13 ;
UPDATE `privilege` SET name='preparation_surat_jalan' WHERE id=14 ;
