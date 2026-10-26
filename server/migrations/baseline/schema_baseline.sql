--
-- PostgreSQL database dump
--



SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

SET default_table_access_method = heap;

--
-- Name: ambient_sounds; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.ambient_sounds (
    id bigint NOT NULL,
    key character varying(50) NOT NULL,
    name character varying(50) DEFAULT ''::character varying NOT NULL,
    emoji character varying(16) DEFAULT ''::character varying NOT NULL,
    url character varying(500) NOT NULL,
    credit character varying(200) DEFAULT ''::character varying NOT NULL,
    sort_order integer DEFAULT 0 NOT NULL,
    is_active boolean DEFAULT true NOT NULL,
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone DEFAULT now()
);


--
-- Name: ambient_sounds_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.ambient_sounds_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: ambient_sounds_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.ambient_sounds_id_seq OWNED BY public.ambient_sounds.id;


--
-- Name: bookmarks; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.bookmarks (
    id bigint NOT NULL,
    user_id bigint NOT NULL,
    post_id bigint NOT NULL,
    created_at timestamp with time zone DEFAULT now()
);


--
-- Name: TABLE bookmarks; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON TABLE public.bookmarks IS '收藏表';


--
-- Name: bookmarks_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.bookmarks_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: bookmarks_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.bookmarks_id_seq OWNED BY public.bookmarks.id;


--
-- Name: comments; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.comments (
    id bigint NOT NULL,
    post_id bigint NOT NULL,
    user_id bigint NOT NULL,
    parent_id bigint,
    content text NOT NULL,
    like_count integer DEFAULT 0,
    status integer DEFAULT 1,
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone DEFAULT now(),
    deleted_at timestamp with time zone,
    is_anonymous boolean DEFAULT false,
    anonymous_name character varying(50)
);


--
-- Name: TABLE comments; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON TABLE public.comments IS '评论表';


--
-- Name: COLUMN comments.parent_id; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.comments.parent_id IS '父评论ID，NULL表示主评论';


--
-- Name: COLUMN comments.status; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.comments.status IS '状态: 0-草稿, 1-已发布, 2-隐藏, 3-删除';


--
-- Name: COLUMN comments.is_anonymous; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.comments.is_anonymous IS '是否匿名评论';


--
-- Name: COLUMN comments.anonymous_name; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.comments.anonymous_name IS '匿名展示名称';


--
-- Name: comments_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.comments_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: comments_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.comments_id_seq OWNED BY public.comments.id;


--
-- Name: daily_checkins; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.daily_checkins (
    id bigint NOT NULL,
    user_id bigint NOT NULL,
    mood character varying(20) NOT NULL,
    content text,
    images text,
    is_public boolean DEFAULT true,
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone DEFAULT now()
);


--
-- Name: daily_checkins_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.daily_checkins_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: daily_checkins_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.daily_checkins_id_seq OWNED BY public.daily_checkins.id;


--
-- Name: daily_images; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.daily_images (
    id bigint NOT NULL,
    url character varying(500) NOT NULL,
    description character varying(200) DEFAULT ''::character varying NOT NULL,
    credit character varying(200) DEFAULT ''::character varying NOT NULL,
    sort_order integer DEFAULT 0 NOT NULL,
    is_active boolean DEFAULT true NOT NULL,
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone DEFAULT now()
);


--
-- Name: daily_images_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.daily_images_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: daily_images_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.daily_images_id_seq OWNED BY public.daily_images.id;


--
-- Name: daily_stats; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.daily_stats (
    id bigint NOT NULL,
    stat_date character varying(10) NOT NULL,
    dimension character varying(20) NOT NULL,
    label character varying(100) NOT NULL,
    count bigint DEFAULT 0 NOT NULL,
    created_at timestamp with time zone DEFAULT now()
);


--
-- Name: daily_stats_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.daily_stats_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: daily_stats_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.daily_stats_id_seq OWNED BY public.daily_stats.id;


--
-- Name: email_verifications; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.email_verifications (
    id bigint NOT NULL,
    user_id bigint NOT NULL,
    token character varying(100) NOT NULL,
    expires_at timestamp with time zone NOT NULL,
    verified_at timestamp with time zone,
    created_at timestamp with time zone DEFAULT now()
);


--
-- Name: email_verifications_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.email_verifications_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: email_verifications_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.email_verifications_id_seq OWNED BY public.email_verifications.id;


--
-- Name: feedbacks; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.feedbacks (
    id bigint NOT NULL,
    user_id bigint NOT NULL,
    type character varying(20) DEFAULT 'other'::character varying NOT NULL,
    title character varying(200) NOT NULL,
    content text NOT NULL,
    contact character varying(100),
    status smallint DEFAULT 0 NOT NULL,
    reply text,
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone DEFAULT now()
);


--
-- Name: feedbacks_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.feedbacks_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: feedbacks_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.feedbacks_id_seq OWNED BY public.feedbacks.id;


--
-- Name: follows; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.follows (
    id bigint NOT NULL,
    follower_id bigint NOT NULL,
    following_id bigint NOT NULL,
    created_at timestamp with time zone DEFAULT now()
);


--
-- Name: TABLE follows; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON TABLE public.follows IS '关注关系表';


--
-- Name: COLUMN follows.follower_id; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.follows.follower_id IS '关注者ID';


--
-- Name: COLUMN follows.following_id; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.follows.following_id IS '被关注者ID';


--
-- Name: follows_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.follows_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: follows_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.follows_id_seq OWNED BY public.follows.id;


--
-- Name: journals; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.journals (
    id bigint NOT NULL,
    user_id bigint NOT NULL,
    content text NOT NULL,
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone DEFAULT now()
);


--
-- Name: journals_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.journals_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: journals_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.journals_id_seq OWNED BY public.journals.id;


--
-- Name: likes; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.likes (
    id bigint NOT NULL,
    user_id bigint NOT NULL,
    post_id bigint NOT NULL,
    created_at timestamp with time zone DEFAULT now()
);


--
-- Name: TABLE likes; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON TABLE public.likes IS '点赞表';


--
-- Name: likes_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.likes_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: likes_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.likes_id_seq OWNED BY public.likes.id;


--
-- Name: notifications; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.notifications (
    id bigint NOT NULL,
    user_id bigint NOT NULL,
    actor_id bigint NOT NULL,
    type character varying(20) NOT NULL,
    target_id bigint,
    target_type character varying(20),
    content text,
    is_read boolean DEFAULT false,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);


--
-- Name: notifications_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.notifications_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: notifications_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.notifications_id_seq OWNED BY public.notifications.id;


--
-- Name: post_tags; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.post_tags (
    post_id integer NOT NULL,
    tag_id integer NOT NULL
);


--
-- Name: posts; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.posts (
    id bigint NOT NULL,
    user_id bigint NOT NULL,
    type character varying(20) NOT NULL,
    title character varying(200) NOT NULL,
    content text NOT NULL,
    images jsonb DEFAULT '[]'::jsonb,
    tags jsonb DEFAULT '[]'::jsonb,
    like_count integer DEFAULT 0,
    comment_count integer DEFAULT 0,
    status integer DEFAULT 1,
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone DEFAULT now(),
    deleted_at timestamp with time zone,
    is_top boolean DEFAULT false,
    is_featured boolean DEFAULT false,
    review_status character varying(20) DEFAULT 'approved'::character varying,
    reviewed_at timestamp with time zone,
    reviewed_by bigint,
    reject_reason text,
    is_anonymous boolean DEFAULT false,
    anonymous_name character varying(50)
);


--
-- Name: TABLE posts; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON TABLE public.posts IS '帖子表';


--
-- Name: COLUMN posts.type; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.posts.type IS '帖子类型: emotion-情绪, insight-认知, case-案例, action-行动';


--
-- Name: COLUMN posts.images; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.posts.images IS '图片URL数组';


--
-- Name: COLUMN posts.tags; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.posts.tags IS '标签数组';


--
-- Name: COLUMN posts.status; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.posts.status IS '状态: 0-草稿, 1-已发布, 2-隐藏, 3-删除';


--
-- Name: COLUMN posts.is_anonymous; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.posts.is_anonymous IS '是否匿名发布';


--
-- Name: COLUMN posts.anonymous_name; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.posts.anonymous_name IS '匿名展示名称';


--
-- Name: posts_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.posts_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: posts_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.posts_id_seq OWNED BY public.posts.id;


--
-- Name: readings; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.readings (
    id bigint NOT NULL,
    source character varying(200) NOT NULL,
    dynasty character varying(50),
    author character varying(100),
    title character varying(200),
    content text NOT NULL,
    interpretation text,
    tags character varying(500),
    image_url character varying(500),
    is_active boolean DEFAULT true NOT NULL,
    created_at timestamp with time zone DEFAULT now()
);


--
-- Name: readings_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.readings_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: readings_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.readings_id_seq OWNED BY public.readings.id;


--
-- Name: reports; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.reports (
    id bigint NOT NULL,
    user_id bigint NOT NULL,
    target_type character varying(20) NOT NULL,
    target_id bigint NOT NULL,
    reason character varying(500) NOT NULL,
    status integer DEFAULT 1,
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone DEFAULT now()
);


--
-- Name: reports_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.reports_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: reports_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.reports_id_seq OWNED BY public.reports.id;


--
-- Name: sensitive_words; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.sensitive_words (
    id bigint NOT NULL,
    word character varying(100) NOT NULL
);


--
-- Name: sensitive_words_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.sensitive_words_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: sensitive_words_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.sensitive_words_id_seq OWNED BY public.sensitive_words.id;


--
-- Name: survival_tips; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.survival_tips (
    id bigint NOT NULL,
    content text NOT NULL,
    category character varying(50) DEFAULT 'general'::character varying,
    is_active boolean DEFAULT true,
    created_at timestamp with time zone DEFAULT now()
);


--
-- Name: survival_tips_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.survival_tips_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: survival_tips_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.survival_tips_id_seq OWNED BY public.survival_tips.id;


--
-- Name: system_settings; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.system_settings (
    key character varying(100) NOT NULL,
    value text DEFAULT ''::text NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL
);


--
-- Name: TABLE system_settings; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON TABLE public.system_settings IS '系统运行时设置，管理员后台可编辑';


--
-- Name: COLUMN system_settings.key; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.system_settings.key IS '设置键，如 smtp.host / smtp.enabled';


--
-- Name: COLUMN system_settings.value; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.system_settings.value IS '设置值，统一以字符串存储，由 service 层解析';


--
-- Name: tags; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.tags (
    id integer NOT NULL,
    name character varying(50) NOT NULL,
    created_at timestamp without time zone DEFAULT now()
);


--
-- Name: tags_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.tags_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: tags_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.tags_id_seq OWNED BY public.tags.id;


--
-- Name: users; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.users (
    id bigint NOT NULL,
    phone character varying(20),
    email character varying(255),
    password character varying(255) NOT NULL,
    nickname character varying(50) NOT NULL,
    company character varying(100) DEFAULT ''::character varying,
    occupation character varying(100) DEFAULT ''::character varying,
    avatar character varying(500) DEFAULT ''::character varying,
    status smallint DEFAULT 1,
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone DEFAULT now(),
    bio character varying(500) DEFAULT ''::character varying,
    following_count integer DEFAULT 0,
    follower_count integer DEFAULT 0,
    is_admin boolean DEFAULT false,
    email_verified boolean DEFAULT false,
    last_work_date timestamp without time zone,
    city character varying(50) DEFAULT ''::character varying,
    industry character varying(100) DEFAULT ''::character varying,
    left_job_at timestamp without time zone
);


--
-- Name: TABLE users; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON TABLE public.users IS '用户表';


--
-- Name: COLUMN users.phone; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.users.phone IS '手机号';


--
-- Name: COLUMN users.email; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.users.email IS '邮箱';


--
-- Name: COLUMN users.password; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.users.password IS '密码哈希';


--
-- Name: COLUMN users.nickname; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.users.nickname IS '昵称';


--
-- Name: COLUMN users.company; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.users.company IS '前公司';


--
-- Name: COLUMN users.occupation; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.users.occupation IS '职业';


--
-- Name: COLUMN users.avatar; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.users.avatar IS '头像URL';


--
-- Name: COLUMN users.status; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.users.status IS '状态: 1-正常, 0-禁用';


--
-- Name: COLUMN users.bio; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.users.bio IS '个人简介';


--
-- Name: COLUMN users.following_count; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.users.following_count IS '关注数';


--
-- Name: COLUMN users.follower_count; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.users.follower_count IS '粉丝数';


--
-- Name: COLUMN users.email_verified; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.users.email_verified IS '邮箱是否已验证';


--
-- Name: COLUMN users.last_work_date; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON COLUMN public.users.last_work_date IS '最后工作日，用于计算自由天数';


--
-- Name: users_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.users_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: users_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.users_id_seq OWNED BY public.users.id;


--
-- Name: ambient_sounds id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.ambient_sounds ALTER COLUMN id SET DEFAULT nextval('public.ambient_sounds_id_seq'::regclass);


--
-- Name: bookmarks id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.bookmarks ALTER COLUMN id SET DEFAULT nextval('public.bookmarks_id_seq'::regclass);


--
-- Name: comments id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.comments ALTER COLUMN id SET DEFAULT nextval('public.comments_id_seq'::regclass);


--
-- Name: daily_checkins id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.daily_checkins ALTER COLUMN id SET DEFAULT nextval('public.daily_checkins_id_seq'::regclass);


--
-- Name: daily_images id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.daily_images ALTER COLUMN id SET DEFAULT nextval('public.daily_images_id_seq'::regclass);


--
-- Name: daily_stats id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.daily_stats ALTER COLUMN id SET DEFAULT nextval('public.daily_stats_id_seq'::regclass);


--
-- Name: email_verifications id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.email_verifications ALTER COLUMN id SET DEFAULT nextval('public.email_verifications_id_seq'::regclass);


--
-- Name: feedbacks id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.feedbacks ALTER COLUMN id SET DEFAULT nextval('public.feedbacks_id_seq'::regclass);


--
-- Name: follows id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.follows ALTER COLUMN id SET DEFAULT nextval('public.follows_id_seq'::regclass);


--
-- Name: journals id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.journals ALTER COLUMN id SET DEFAULT nextval('public.journals_id_seq'::regclass);


--
-- Name: likes id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.likes ALTER COLUMN id SET DEFAULT nextval('public.likes_id_seq'::regclass);


--
-- Name: notifications id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.notifications ALTER COLUMN id SET DEFAULT nextval('public.notifications_id_seq'::regclass);


--
-- Name: posts id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.posts ALTER COLUMN id SET DEFAULT nextval('public.posts_id_seq'::regclass);


--
-- Name: readings id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.readings ALTER COLUMN id SET DEFAULT nextval('public.readings_id_seq'::regclass);


--
-- Name: reports id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.reports ALTER COLUMN id SET DEFAULT nextval('public.reports_id_seq'::regclass);


--
-- Name: sensitive_words id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.sensitive_words ALTER COLUMN id SET DEFAULT nextval('public.sensitive_words_id_seq'::regclass);


--
-- Name: survival_tips id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.survival_tips ALTER COLUMN id SET DEFAULT nextval('public.survival_tips_id_seq'::regclass);


--
-- Name: tags id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.tags ALTER COLUMN id SET DEFAULT nextval('public.tags_id_seq'::regclass);


--
-- Name: users id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.users ALTER COLUMN id SET DEFAULT nextval('public.users_id_seq'::regclass);


--
-- Name: ambient_sounds ambient_sounds_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.ambient_sounds
    ADD CONSTRAINT ambient_sounds_pkey PRIMARY KEY (id);


--
-- Name: bookmarks bookmarks_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.bookmarks
    ADD CONSTRAINT bookmarks_pkey PRIMARY KEY (id);


--
-- Name: bookmarks bookmarks_user_id_post_id_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.bookmarks
    ADD CONSTRAINT bookmarks_user_id_post_id_key UNIQUE (user_id, post_id);


--
-- Name: comments comments_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.comments
    ADD CONSTRAINT comments_pkey PRIMARY KEY (id);


--
-- Name: daily_checkins daily_checkins_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.daily_checkins
    ADD CONSTRAINT daily_checkins_pkey PRIMARY KEY (id);


--
-- Name: daily_images daily_images_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.daily_images
    ADD CONSTRAINT daily_images_pkey PRIMARY KEY (id);


--
-- Name: daily_stats daily_stats_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.daily_stats
    ADD CONSTRAINT daily_stats_pkey PRIMARY KEY (id);


--
-- Name: email_verifications email_verifications_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.email_verifications
    ADD CONSTRAINT email_verifications_pkey PRIMARY KEY (id);


--
-- Name: email_verifications email_verifications_token_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.email_verifications
    ADD CONSTRAINT email_verifications_token_key UNIQUE (token);


--
-- Name: feedbacks feedbacks_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.feedbacks
    ADD CONSTRAINT feedbacks_pkey PRIMARY KEY (id);


--
-- Name: follows follows_follower_id_following_id_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.follows
    ADD CONSTRAINT follows_follower_id_following_id_key UNIQUE (follower_id, following_id);


--
-- Name: follows follows_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.follows
    ADD CONSTRAINT follows_pkey PRIMARY KEY (id);


--
-- Name: journals journals_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.journals
    ADD CONSTRAINT journals_pkey PRIMARY KEY (id);


--
-- Name: likes likes_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.likes
    ADD CONSTRAINT likes_pkey PRIMARY KEY (id);


--
-- Name: likes likes_user_id_post_id_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.likes
    ADD CONSTRAINT likes_user_id_post_id_key UNIQUE (user_id, post_id);


--
-- Name: notifications notifications_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.notifications
    ADD CONSTRAINT notifications_pkey PRIMARY KEY (id);


--
-- Name: post_tags post_tags_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.post_tags
    ADD CONSTRAINT post_tags_pkey PRIMARY KEY (post_id, tag_id);


--
-- Name: posts posts_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.posts
    ADD CONSTRAINT posts_pkey PRIMARY KEY (id);


--
-- Name: readings readings_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.readings
    ADD CONSTRAINT readings_pkey PRIMARY KEY (id);


--
-- Name: reports reports_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.reports
    ADD CONSTRAINT reports_pkey PRIMARY KEY (id);


--
-- Name: sensitive_words sensitive_words_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.sensitive_words
    ADD CONSTRAINT sensitive_words_pkey PRIMARY KEY (id);


--
-- Name: sensitive_words sensitive_words_word_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.sensitive_words
    ADD CONSTRAINT sensitive_words_word_key UNIQUE (word);


--
-- Name: survival_tips survival_tips_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.survival_tips
    ADD CONSTRAINT survival_tips_pkey PRIMARY KEY (id);


--
-- Name: system_settings system_settings_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.system_settings
    ADD CONSTRAINT system_settings_pkey PRIMARY KEY (key);


--
-- Name: tags tags_name_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.tags
    ADD CONSTRAINT tags_name_key UNIQUE (name);


--
-- Name: tags tags_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.tags
    ADD CONSTRAINT tags_pkey PRIMARY KEY (id);


--
-- Name: users users_email_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_email_key UNIQUE (email);


--
-- Name: users users_phone_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_phone_key UNIQUE (phone);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- Name: idx_ambient_sounds_is_active; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_ambient_sounds_is_active ON public.ambient_sounds USING btree (is_active);


--
-- Name: idx_ambient_sounds_sort_order; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_ambient_sounds_sort_order ON public.ambient_sounds USING btree (sort_order);


--
-- Name: idx_bookmarks_post_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_bookmarks_post_id ON public.bookmarks USING btree (post_id);


--
-- Name: idx_bookmarks_user_created; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_bookmarks_user_created ON public.bookmarks USING btree (user_id, created_at DESC);


--
-- Name: INDEX idx_bookmarks_user_created; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON INDEX public.idx_bookmarks_user_created IS '收藏列表索引：优化 ListBookmarks 分页';


--
-- Name: idx_bookmarks_user_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_bookmarks_user_id ON public.bookmarks USING btree (user_id);


--
-- Name: idx_bookmarks_user_post; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_bookmarks_user_post ON public.bookmarks USING btree (user_id, post_id);


--
-- Name: INDEX idx_bookmarks_user_post; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON INDEX public.idx_bookmarks_user_post IS '收藏精确查询复合索引：优化 IsBookmarked/Unbookmark';


--
-- Name: idx_comments_created_at; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_comments_created_at ON public.comments USING btree (created_at);


--
-- Name: idx_comments_parent_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_comments_parent_id ON public.comments USING btree (parent_id) WHERE (parent_id IS NOT NULL);


--
-- Name: idx_comments_post_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_comments_post_id ON public.comments USING btree (post_id);


--
-- Name: idx_comments_post_status_created; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_comments_post_status_created ON public.comments USING btree (post_id, status, created_at);


--
-- Name: INDEX idx_comments_post_status_created; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON INDEX public.idx_comments_post_status_created IS '评论列表查询复合索引';


--
-- Name: idx_comments_status; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_comments_status ON public.comments USING btree (status);


--
-- Name: idx_comments_user_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_comments_user_id ON public.comments USING btree (user_id);


--
-- Name: idx_daily_checkins_created_at; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_daily_checkins_created_at ON public.daily_checkins USING btree (created_at DESC);


--
-- Name: idx_daily_checkins_public_created; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_daily_checkins_public_created ON public.daily_checkins USING btree (is_public, created_at DESC);


--
-- Name: idx_daily_checkins_user_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_daily_checkins_user_id ON public.daily_checkins USING btree (user_id);


--
-- Name: idx_daily_images_is_active; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_daily_images_is_active ON public.daily_images USING btree (is_active);


--
-- Name: idx_daily_images_sort_order; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_daily_images_sort_order ON public.daily_images USING btree (sort_order);


--
-- Name: idx_daily_stats_date; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_daily_stats_date ON public.daily_stats USING btree (stat_date);


--
-- Name: idx_daily_stats_date_dim; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_daily_stats_date_dim ON public.daily_stats USING btree (stat_date, dimension);


--
-- Name: idx_daily_stats_dimension; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_daily_stats_dimension ON public.daily_stats USING btree (dimension);


--
-- Name: idx_daily_stats_label; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_daily_stats_label ON public.daily_stats USING btree (label);


--
-- Name: idx_email_verifications_user_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_email_verifications_user_id ON public.email_verifications USING btree (user_id);


--
-- Name: idx_feedbacks_created_at; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_feedbacks_created_at ON public.feedbacks USING btree (created_at DESC);


--
-- Name: idx_feedbacks_status; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_feedbacks_status ON public.feedbacks USING btree (status);


--
-- Name: idx_feedbacks_user_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_feedbacks_user_id ON public.feedbacks USING btree (user_id);


--
-- Name: idx_follows_follower; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_follows_follower ON public.follows USING btree (follower_id);


--
-- Name: idx_follows_follower_following; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_follows_follower_following ON public.follows USING btree (follower_id, following_id);


--
-- Name: INDEX idx_follows_follower_following; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON INDEX public.idx_follows_follower_following IS '关注关系精确查询索引：优化 IsFollowing/Unfollow';


--
-- Name: idx_follows_follower_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_follows_follower_id ON public.follows USING btree (follower_id);


--
-- Name: idx_follows_following; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_follows_following ON public.follows USING btree (following_id);


--
-- Name: idx_follows_following_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_follows_following_id ON public.follows USING btree (following_id);


--
-- Name: idx_journals_created_at; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_journals_created_at ON public.journals USING btree (created_at DESC);


--
-- Name: idx_journals_user_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_journals_user_id ON public.journals USING btree (user_id);


--
-- Name: idx_likes_post_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_likes_post_id ON public.likes USING btree (post_id);


--
-- Name: idx_likes_user_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_likes_user_id ON public.likes USING btree (user_id);


--
-- Name: idx_likes_user_post; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_likes_user_post ON public.likes USING btree (user_id, post_id);


--
-- Name: INDEX idx_likes_user_post; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON INDEX public.idx_likes_user_post IS '点赞查询复合索引：优化 IsLiked/Unlike';


--
-- Name: idx_notifications_created_at; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_notifications_created_at ON public.notifications USING btree (created_at DESC);


--
-- Name: idx_notifications_type; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_notifications_type ON public.notifications USING btree (type);


--
-- Name: idx_notifications_user_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_notifications_user_id ON public.notifications USING btree (user_id);


--
-- Name: idx_notifications_user_read; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_notifications_user_read ON public.notifications USING btree (user_id, is_read);


--
-- Name: idx_notifications_user_read_created; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_notifications_user_read_created ON public.notifications USING btree (user_id, is_read, created_at DESC);


--
-- Name: INDEX idx_notifications_user_read_created; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON INDEX public.idx_notifications_user_read_created IS '通知列表查询复合索引';


--
-- Name: idx_post_tags_post_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_post_tags_post_id ON public.post_tags USING btree (post_id);


--
-- Name: idx_post_tags_tag_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_post_tags_tag_id ON public.post_tags USING btree (tag_id);


--
-- Name: idx_posts_created_at; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_posts_created_at ON public.posts USING btree (created_at DESC);


--
-- Name: idx_posts_deleted_at; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_posts_deleted_at ON public.posts USING btree (deleted_at) WHERE (deleted_at IS NOT NULL);


--
-- Name: idx_posts_is_featured; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_posts_is_featured ON public.posts USING btree (is_featured);


--
-- Name: idx_posts_is_top; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_posts_is_top ON public.posts USING btree (is_top);


--
-- Name: idx_posts_review_status; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_posts_review_status ON public.posts USING btree (review_status);


--
-- Name: idx_posts_status; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_posts_status ON public.posts USING btree (status);


--
-- Name: idx_posts_status_top_featured_created; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_posts_status_top_featured_created ON public.posts USING btree (status, is_top, is_featured, created_at DESC);


--
-- Name: idx_posts_status_type_created; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_posts_status_type_created ON public.posts USING btree (status, type, created_at DESC);


--
-- Name: INDEX idx_posts_status_type_created; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON INDEX public.idx_posts_status_type_created IS '帖子列表查询复合索引';


--
-- Name: idx_posts_type; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_posts_type ON public.posts USING btree (type);


--
-- Name: idx_posts_user_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_posts_user_id ON public.posts USING btree (user_id);


--
-- Name: idx_posts_user_status_created; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_posts_user_status_created ON public.posts USING btree (user_id, status, created_at DESC);


--
-- Name: INDEX idx_posts_user_status_created; Type: COMMENT; Schema: public; Owner: -
--

COMMENT ON INDEX public.idx_posts_user_status_created IS '用户帖子查询复合索引';


--
-- Name: idx_readings_is_active; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_readings_is_active ON public.readings USING btree (is_active);


--
-- Name: idx_reports_status; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_reports_status ON public.reports USING btree (status);


--
-- Name: idx_reports_target; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_reports_target ON public.reports USING btree (target_type, target_id);


--
-- Name: idx_reports_user_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_reports_user_id ON public.reports USING btree (user_id);


--
-- Name: idx_survival_tips_active; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_survival_tips_active ON public.survival_tips USING btree (is_active);


--
-- Name: idx_tags_name; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_tags_name ON public.tags USING btree (name);


--
-- Name: idx_users_city; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_users_city ON public.users USING btree (city);


--
-- Name: idx_users_city_industry_status; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_users_city_industry_status ON public.users USING btree (city, industry, status);


--
-- Name: idx_users_email; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_users_email ON public.users USING btree (email) WHERE (email IS NOT NULL);


--
-- Name: idx_users_industry; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_users_industry ON public.users USING btree (industry);


--
-- Name: idx_users_nickname; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_users_nickname ON public.users USING btree (nickname);


--
-- Name: idx_users_phone; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_users_phone ON public.users USING btree (phone) WHERE (phone IS NOT NULL);


--
-- Name: idx_users_status; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_users_status ON public.users USING btree (status);


--
-- Name: daily_checkins daily_checkins_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.daily_checkins
    ADD CONSTRAINT daily_checkins_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id);


--
-- Name: email_verifications email_verifications_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.email_verifications
    ADD CONSTRAINT email_verifications_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id);


--
-- Name: bookmarks fk_bookmarks_post_id; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.bookmarks
    ADD CONSTRAINT fk_bookmarks_post_id FOREIGN KEY (post_id) REFERENCES public.posts(id);


--
-- Name: bookmarks fk_bookmarks_user_id; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.bookmarks
    ADD CONSTRAINT fk_bookmarks_user_id FOREIGN KEY (user_id) REFERENCES public.users(id);


--
-- Name: comments fk_comments_parent_id; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.comments
    ADD CONSTRAINT fk_comments_parent_id FOREIGN KEY (parent_id) REFERENCES public.comments(id);


--
-- Name: comments fk_comments_post_id; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.comments
    ADD CONSTRAINT fk_comments_post_id FOREIGN KEY (post_id) REFERENCES public.posts(id);


--
-- Name: comments fk_comments_user_id; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.comments
    ADD CONSTRAINT fk_comments_user_id FOREIGN KEY (user_id) REFERENCES public.users(id);


--
-- Name: follows fk_follows_follower_id; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.follows
    ADD CONSTRAINT fk_follows_follower_id FOREIGN KEY (follower_id) REFERENCES public.users(id);


--
-- Name: follows fk_follows_following_id; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.follows
    ADD CONSTRAINT fk_follows_following_id FOREIGN KEY (following_id) REFERENCES public.users(id);


--
-- Name: likes fk_likes_post_id; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.likes
    ADD CONSTRAINT fk_likes_post_id FOREIGN KEY (post_id) REFERENCES public.posts(id);


--
-- Name: likes fk_likes_user_id; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.likes
    ADD CONSTRAINT fk_likes_user_id FOREIGN KEY (user_id) REFERENCES public.users(id);


--
-- Name: posts fk_posts_user_id; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.posts
    ADD CONSTRAINT fk_posts_user_id FOREIGN KEY (user_id) REFERENCES public.users(id);


--
-- Name: notifications notifications_actor_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.notifications
    ADD CONSTRAINT notifications_actor_id_fkey FOREIGN KEY (actor_id) REFERENCES public.users(id);


--
-- Name: notifications notifications_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.notifications
    ADD CONSTRAINT notifications_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id);


--
-- Name: post_tags post_tags_post_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.post_tags
    ADD CONSTRAINT post_tags_post_id_fkey FOREIGN KEY (post_id) REFERENCES public.posts(id) ON DELETE CASCADE;


--
-- Name: post_tags post_tags_tag_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.post_tags
    ADD CONSTRAINT post_tags_tag_id_fkey FOREIGN KEY (tag_id) REFERENCES public.tags(id) ON DELETE CASCADE;


--
-- PostgreSQL database dump complete
--


