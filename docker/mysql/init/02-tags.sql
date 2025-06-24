-- 标签索引表
CREATE TABLE IF NOT EXISTS `aq3cms_tagindex` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `tag` varchar(30) NOT NULL DEFAULT '',
  `count` int(11) NOT NULL DEFAULT '0',
  `rank` int(11) NOT NULL DEFAULT '0',
  `ishot` tinyint(1) NOT NULL DEFAULT '0',
  `addtime` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `lastuse` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `tagpinyin` varchar(100) NOT NULL DEFAULT '',
  PRIMARY KEY (`id`),
  UNIQUE KEY `tag` (`tag`),
  KEY `rank` (`rank`),
  KEY `ishot` (`ishot`),
  KEY `count` (`count`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 标签关联表
CREATE TABLE IF NOT EXISTS `aq3cms_taglist` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `aid` int(11) NOT NULL DEFAULT '0',
  `tag` varchar(30) NOT NULL DEFAULT '',
  `addtime` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `aid` (`aid`),
  KEY `tag` (`tag`),
  UNIQUE KEY `aid_tag` (`aid`, `tag`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 插入一些示例标签
INSERT INTO `aq3cms_tagindex` (`tag`, `count`, `rank`, `ishot`, `tagpinyin`) VALUES
('技术', 5, 10, 1, 'jishu'),
('新闻', 8, 20, 1, 'xinwen'),
('产品', 3, 15, 0, 'chanpin'),
('教程', 6, 12, 1, 'jiaocheng'),
('开发', 4, 8, 0, 'kaifa');
