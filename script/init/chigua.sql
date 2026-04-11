-- 吃瓜网数据库脚本

DROP TABLE IF EXISTS users;
CREATE TABLE users (
    id INT8 GENERATED ALWAYS AS IDENTITY,
    username VARCHAR(50) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    nickname VARCHAR(50),
    role VARCHAR(20) DEFAULT 'user',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    update_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id)
);
COMMENT ON COLUMN users.id IS 'ID';
COMMENT ON COLUMN users.username IS '用户名';
COMMENT ON COLUMN users.password IS '密码';
COMMENT ON COLUMN users.nickname IS '昵称';
COMMENT ON COLUMN users.role IS '角色';
COMMENT ON COLUMN users.created_at IS '创建时间';
COMMENT ON COLUMN users.update_at IS '更新时间';
COMMENT ON TABLE users IS '用户表';

INSERT INTO users (username, password, nickname, role) VALUES ('admin', '8d969eef6ecad3c29a3a629280e686cf0c3f5d5a86aff3ca12020c923adc6c92', 'sylvie233', 'admin');

DROP TABLE IF EXISTS article;
CREATE TABLE article (
    id INT8 GENERATED ALWAYS AS IDENTITY,
    author_id INT8,
    category_id INT8,
    title VARCHAR(255) NOT NULL,
    content TEXT,
    cover_image VARCHAR(255),
    status INT2,
    publish_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    update_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id)
);
COMMENT ON COLUMN article.id IS 'ID';
COMMENT ON COLUMN article.author_id IS '作者ID';
COMMENT ON COLUMN article.category_id IS '分类ID';
COMMENT ON COLUMN article.title IS '标题';
COMMENT ON COLUMN article.content IS '内容';
COMMENT ON COLUMN article.cover_image IS '首页图';
COMMENT ON COLUMN article.status IS '状态;0：未发布、1：已发布';
COMMENT ON COLUMN article.publish_at IS '发布时间';
COMMENT ON COLUMN article.created_at IS '创建时间';
COMMENT ON COLUMN article.update_at IS '更新时间';
COMMENT ON TABLE article IS '文章表';

DROP TABLE IF EXISTS comment;
CREATE TABLE comment (
    id INT8 GENERATED ALWAYS AS IDENTITY,
    parant_id INT8 NOT NULL DEFAULT 0,
    article_id INT8 NOT NULL,
    reply_user_id INT8 DEFAULT 0,
    user_id INT8 NOT NULL,
    content TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id)
);
COMMENT ON COLUMN comment.id IS 'ID';
COMMENT ON COLUMN comment.parant_id IS '父级评论ID';
COMMENT ON COLUMN comment.article_id IS '文章ID';
COMMENT ON COLUMN comment.reply_user_id IS '回复用户ID;可空';
COMMENT ON COLUMN comment.user_id IS '用户ID';
COMMENT ON COLUMN comment.content IS '评论内容';
COMMENT ON COLUMN comment.created_at IS '创建时间';
COMMENT ON TABLE comment IS '评论表';


DROP TABLE IF EXISTS tag;
CREATE TABLE tag (
    id INT8 GENERATED ALWAYS AS IDENTITY,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    update_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id)
);
COMMENT ON COLUMN tag.id IS 'ID';
COMMENT ON COLUMN tag.name IS '标签名';
COMMENT ON COLUMN tag.created_at IS '创建时间';
COMMENT ON COLUMN tag.update_at IS '更新时间';
COMMENT ON TABLE tag IS '标签表';

DROP TABLE IF EXISTS article_tag;
CREATE TABLE article_tag (
    article_id INT8,
    tag_id INT8,
    PRIMARY KEY (article_id,tag_id)
);
COMMENT ON COLUMN article_tag.article_id IS '文章ID';
COMMENT ON COLUMN article_tag.tag_id IS '标签ID';
COMMENT ON TABLE article_tag IS '文章标签关联表';

DROP TABLE IF EXISTS category;
CREATE TABLE category (
    id INT8 GENERATED ALWAYS AS IDENTITY,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    update_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id)
);
COMMENT ON COLUMN category.id IS 'ID';
COMMENT ON COLUMN category.name IS '分类名';
COMMENT ON COLUMN category.created_at IS '创建时间';
COMMENT ON COLUMN category.update_at IS '更新时间';
COMMENT ON TABLE category IS '分类表';
