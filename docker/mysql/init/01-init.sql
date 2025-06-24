-- aq3cms 数据库初始化脚本

-- 创建数据库
CREATE DATABASE IF NOT EXISTS `aq3cms` DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

USE `aq3cms`;

-- 管理员表
CREATE TABLE IF NOT EXISTS `aq3cms_admin` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `username` varchar(50) NOT NULL DEFAULT '',
  `password` varchar(255) NOT NULL DEFAULT '',
  `rank` int(11) NOT NULL DEFAULT '10',
  `email` varchar(100) NOT NULL DEFAULT '',
  `status` tinyint(1) NOT NULL DEFAULT '1',
  `last_login` datetime DEFAULT NULL,
  `last_ip` varchar(45) DEFAULT '',
  `login_count` int(11) NOT NULL DEFAULT '0',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `username` (`username`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 栏目表
CREATE TABLE IF NOT EXISTS `aq3cms_arctype` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `reid` int(11) NOT NULL DEFAULT '0',
  `topid` int(11) NOT NULL DEFAULT '0',
  `sortrank` int(11) NOT NULL DEFAULT '50',
  `typename` varchar(60) NOT NULL DEFAULT '',
  `typedir` varchar(60) NOT NULL DEFAULT '',
  `isdefault` tinyint(1) NOT NULL DEFAULT '0',
  `defaultname` varchar(60) NOT NULL DEFAULT 'index.html',
  `issend` tinyint(1) NOT NULL DEFAULT '1',
  `channeltype` int(11) NOT NULL DEFAULT '1',
  `maxpage` int(11) NOT NULL DEFAULT '-1',
  `ispart` tinyint(1) NOT NULL DEFAULT '0',
  `corank` int(11) NOT NULL DEFAULT '0',
  `tempindex` varchar(60) NOT NULL DEFAULT '',
  `templist` varchar(60) NOT NULL DEFAULT '',
  `temparticle` varchar(60) NOT NULL DEFAULT '',
  `namerule` varchar(60) NOT NULL DEFAULT '',
  `namerule2` varchar(60) NOT NULL DEFAULT '',
  `modname` varchar(20) NOT NULL DEFAULT '',
  `description` text,
  `keywords` varchar(60) NOT NULL DEFAULT '',
  `seotitle` varchar(80) NOT NULL DEFAULT '',
  `moresite` tinyint(1) NOT NULL DEFAULT '0',
  `sitepath` varchar(60) NOT NULL DEFAULT '',
  `siteurl` varchar(200) NOT NULL DEFAULT '',
  `ishidden` tinyint(1) NOT NULL DEFAULT '0',
  `cross` tinyint(1) NOT NULL DEFAULT '0',
  `crossid` text,
  `content` text,
  `smalltypes` text,
  PRIMARY KEY (`id`),
  KEY `reid` (`reid`),
  KEY `topid` (`topid`),
  KEY `sortrank` (`sortrank`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 文档主表
CREATE TABLE IF NOT EXISTS `aq3cms_archives` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `typeid` int(11) NOT NULL DEFAULT '0',
  `typeid2` varchar(90) NOT NULL DEFAULT '0',
  `sortrank` int(11) NOT NULL DEFAULT '0',
  `flag` set('c','h','p','f','s','j','a','b') DEFAULT NULL,
  `ismake` tinyint(1) NOT NULL DEFAULT '0',
  `channel` int(11) NOT NULL DEFAULT '1',
  `arcrank` int(11) NOT NULL DEFAULT '0',
  `click` int(11) NOT NULL DEFAULT '0',
  `money` int(11) NOT NULL DEFAULT '0',
  `title` varchar(60) NOT NULL DEFAULT '',
  `shorttitle` varchar(36) NOT NULL DEFAULT '',
  `color` varchar(7) NOT NULL DEFAULT '',
  `writer` varchar(30) NOT NULL DEFAULT '',
  `source` varchar(30) NOT NULL DEFAULT '',
  `litpic` varchar(60) NOT NULL DEFAULT '',
  `pubdate` int(11) NOT NULL DEFAULT '0',
  `senddate` int(11) NOT NULL DEFAULT '0',
  `mid` int(11) NOT NULL DEFAULT '0',
  `keywords` varchar(60) NOT NULL DEFAULT '',
  `lastpost` int(11) NOT NULL DEFAULT '0',
  `scores` int(11) NOT NULL DEFAULT '0',
  `goodpost` int(11) NOT NULL DEFAULT '0',
  `badpost` int(11) NOT NULL DEFAULT '0',
  `voteid` int(11) NOT NULL DEFAULT '0',
  `notpost` tinyint(1) NOT NULL DEFAULT '0',
  `description` varchar(250) NOT NULL DEFAULT '',
  `filename` varchar(60) NOT NULL DEFAULT '',
  `dutyadmin` int(11) NOT NULL DEFAULT '0',
  `tackid` int(11) NOT NULL DEFAULT '0',
  `mtype` int(11) NOT NULL DEFAULT '0',
  `weight` int(11) NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`),
  KEY `typeid` (`typeid`),
  KEY `sortrank` (`sortrank`),
  KEY `pubdate` (`pubdate`),
  KEY `click` (`click`),
  KEY `arcrank` (`arcrank`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 文档附加表
CREATE TABLE IF NOT EXISTS `aq3cms_addonarticle` (
  `aid` int(11) NOT NULL DEFAULT '0',
  `typeid` int(11) NOT NULL DEFAULT '0',
  `body` longtext,
  `redirecturl` varchar(255) NOT NULL DEFAULT '',
  `templet` varchar(30) NOT NULL DEFAULT '',
  `userip` varchar(15) NOT NULL DEFAULT '',
  PRIMARY KEY (`aid`),
  KEY `typeid` (`typeid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 会员表
CREATE TABLE IF NOT EXISTS `aq3cms_member` (
  `mid` int(11) NOT NULL AUTO_INCREMENT,
  `mtype` varchar(20) NOT NULL DEFAULT '',
  `userid` varchar(20) NOT NULL DEFAULT '',
  `pwd` varchar(32) NOT NULL DEFAULT '',
  `uname` varchar(36) NOT NULL DEFAULT '',
  `sex` enum('男','女','保密') NOT NULL DEFAULT '保密',
  `rank` int(11) NOT NULL DEFAULT '10',
  `money` decimal(10,2) NOT NULL DEFAULT '0.00',
  `email` varchar(50) NOT NULL DEFAULT '',
  `scores` int(11) NOT NULL DEFAULT '0',
  `matt` tinyint(1) NOT NULL DEFAULT '0',
  `spacesta` tinyint(1) NOT NULL DEFAULT '0',
  `face` varchar(50) NOT NULL DEFAULT '',
  `safequestion` tinyint(1) NOT NULL DEFAULT '0',
  `safeanswer` varchar(30) NOT NULL DEFAULT '',
  `jointime` int(11) NOT NULL DEFAULT '0',
  `joinip` varchar(16) NOT NULL DEFAULT '',
  `logintime` int(11) NOT NULL DEFAULT '0',
  `loginip` varchar(16) NOT NULL DEFAULT '',
  PRIMARY KEY (`mid`),
  KEY `userid` (`userid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 系统配置表
CREATE TABLE IF NOT EXISTS `aq3cms_sysconfig` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `varname` varchar(20) NOT NULL DEFAULT '',
  `info` varchar(100) NOT NULL DEFAULT '',
  `groupid` tinyint(1) NOT NULL DEFAULT '1',
  `type` varchar(10) NOT NULL DEFAULT 'string',
  `value` text,
  PRIMARY KEY (`id`),
  KEY `varname` (`varname`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 插入默认管理员账号 (用户名: admin, 密码: admin123)
INSERT INTO `aq3cms_admin` (`username`, `password`, `rank`, `email`, `status`) VALUES
('admin', '$argon2id$v=19$m=65536,t=3,p=4$randomsalt$hashedpassword', 10, 'admin@example.com', 1);

-- 插入默认栏目
INSERT INTO `aq3cms_arctype` (`id`, `reid`, `topid`, `sortrank`, `typename`, `typedir`, `isdefault`, `defaultname`, `issend`, `channeltype`, `maxpage`, `ispart`, `corank`, `tempindex`, `templist`, `temparticle`, `namerule`, `namerule2`, `modname`, `description`, `keywords`, `seotitle`) VALUES
(1, 0, 0, 50, '网站首页', '{cmspath}', 1, 'index.html', 1, 1, -1, 0, 0, 'index.htm', 'list.htm', 'article.htm', '', '', '', '网站首页', '', ''),
(2, 0, 0, 50, '新闻中心', 'news', 0, 'index.html', 1, 1, -1, 0, 0, '', 'list_article.htm', 'article_article.htm', '/news/{Y}/{M}{D}/{aid}.html', '/news/list_{tid}_{page}.html', '', '新闻中心', '新闻', '新闻中心'),
(3, 0, 0, 50, '产品展示', 'product', 0, 'index.html', 1, 2, -1, 0, 0, '', 'list_product.htm', 'article_product.htm', '/product/{Y}/{M}{D}/{aid}.html', '/product/list_{tid}_{page}.html', '', '产品展示', '产品', '产品展示'),
(4, 0, 0, 50, '关于我们', 'about', 0, 'index.html', 1, 1, -1, 1, 0, '', '', 'article_article.htm', '', '', '', '关于我们', '关于', '关于我们');

-- 插入系统配置
INSERT INTO `aq3cms_sysconfig` (`varname`, `info`, `groupid`, `type`, `value`) VALUES
('cfg_webname', '网站名称', 1, 'string', 'aq3cms内容管理系统'),
('cfg_weburl', '网站地址', 1, 'string', 'http://localhost'),
('cfg_keywords', '网站关键词', 1, 'string', 'aq3cms,内容管理系统,Go语言'),
('cfg_description', '网站描述', 1, 'string', 'aq3cms是基于Go语言开发的现代化内容管理系统'),
('cfg_icp', 'ICP备案号', 1, 'string', ''),
('cfg_statcode', '统计代码', 1, 'text', ''),
('cfg_copyright', '版权信息', 1, 'string', 'Copyright © 2024 aq3cms. All rights reserved.');
