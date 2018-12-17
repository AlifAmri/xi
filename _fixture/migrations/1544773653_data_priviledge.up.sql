SET FOREIGN_KEY_CHECKS = 0;
INSERT INTO `privilege` (`id`,`name`,`action`,`is_active`,`note`) VALUES
('36', 'preparation', 'qc', '1', 'Dapat melakukan check list qc pada preparation'),
('37', 'preparation', 'checkout', '1', 'Dapat mengeluarkan barang pada preparation'),
('38', 'preparation', 'finish', '1', 'Dapat menyelesaikan preparation');

UPDATE `privilege` SET note='Dapat membuat dan mengelola data preparation' WHERE id=9 ;
UPDATE `privilege` SET action ='publish', note='Dapat melakukan proses publish preparation' WHERE id=10 ;
