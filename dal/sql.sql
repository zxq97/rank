CREATE TABLE `rank_score`
(
    `id`      BIGINT    NOT NULL AUTO_INCREMENT,
    `rank_id` INT       NOT NULL DEFAULT 0 COMMENT '排行榜id',
    `uid`     BIGINT    NOT NULL DEFAULT 0 COMMENT '用户id',
    `score`   INT       NOT NULL DEFAULT 0 COMMENT '分数',
    `rank`    INT       NOT NULL DEFAULT 0 COMMENT '排名，结束后同步',
    `state`   TINYINT   NOT NULL DEFAULT 0 COMMENT '标记位',
    `ctime`   TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `mtime`   TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    UNIQUE uniq_rankid_uid (`rank_id`, `uid`)
)ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COMMENT '榜单分数';

CREATE TABLE `rank_score_tx`
(
    `id`          BIGINT      NOT NULL AUTO_INCREMENT,
    `tx_id`       VARCHAR(64) NOT NULL DEFAULT '' COMMENT '流水编号如:uuid',
    `rank_id`     INT         NOT NULL DEFAULT 0 COMMENT '排行榜id',
    `uid`         BIGINT      NOT NULL DEFAULT 0 COMMENT '用户id',
    `delta_score` INT         NOT NULL DEFAULT 0 COMMENT '单词变化分数',
    `ctime`       TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `mtime`       TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    UNIQUE uniq_txid (`tx_id`),
    KEY           idx_rankid_uid_deltascore (`rank_id`, `uid`, `delta_score`)
)ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COMMENT '榜单分数流水';