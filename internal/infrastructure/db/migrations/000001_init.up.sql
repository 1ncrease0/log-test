CREATE TABLE logs (
    id          BIGSERIAL   PRIMARY KEY,
    path        TEXT        NOT NULL UNIQUE,
    status      TEXT        NOT NULL,
    node_count  INT         NOT NULL DEFAULT 0,
    port_count  INT         NOT NULL DEFAULT 0,
    uploaded_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE nodes (
    id                BIGSERIAL PRIMARY KEY,
    log_id            BIGINT    NOT NULL REFERENCES logs (id) ON DELETE CASCADE,
    node_guid         TEXT      NOT NULL,
    node_desc         TEXT      NOT NULL,
    node_type         INT       NOT NULL,
    num_ports         INT       NOT NULL,
    class_version     INT       NOT NULL,
    base_version      INT       NOT NULL,
    system_image_guid TEXT      NOT NULL,
    port_guid         TEXT      NOT NULL,
    UNIQUE (log_id, node_guid)
);

CREATE TABLE ports (
    id                                      BIGSERIAL PRIMARY KEY,
    log_id                                  BIGINT NOT NULL REFERENCES logs (id)  ON DELETE CASCADE,
    node_id                                 BIGINT NOT NULL REFERENCES nodes (id) ON DELETE CASCADE,
    port_guid                               TEXT   NOT NULL,
    port_num                                INT    NOT NULL,
    m_key                                   TEXT   NOT NULL,
    gid_prfx                                TEXT   NOT NULL,
    msm_lid                                 INT    NOT NULL,
    lid                                     INT    NOT NULL,
    cap_msk                                 BIGINT NOT NULL,
    m_key_lease_period                      INT    NOT NULL,
    diag_code                               INT    NOT NULL,
    link_width_actv                         INT    NOT NULL,
    link_width_sup                          INT    NOT NULL,
    link_width_en                           INT    NOT NULL,
    local_port_num                          INT    NOT NULL,
    link_speed_en                           INT    NOT NULL,
    link_speed_actv                         INT    NOT NULL,
    lmc                                     INT    NOT NULL,
    m_key_prot_bits                         INT    NOT NULL,
    link_down_def_state                     INT    NOT NULL,
    port_phy_state                          INT    NOT NULL,
    port_state                              INT    NOT NULL,
    link_speed_sup                          INT    NOT NULL,
    vl_arb_high_cap                         INT    NOT NULL,
    vl_high_limit                           INT    NOT NULL,
    init_type                               INT    NOT NULL,
    vl_cap                                  INT    NOT NULL,
    msmsl                                   INT    NOT NULL,
    nmtu                                    INT    NOT NULL,
    filter_raw_outb                         INT    NOT NULL,
    filter_raw_inb                          INT    NOT NULL,
    part_enf_outb                           INT    NOT NULL,
    part_enf_inb                            INT    NOT NULL,
    op_vls                                  INT    NOT NULL,
    hoq_life                                INT    NOT NULL,
    vl_stall_cnt                            INT    NOT NULL,
    mtu_cap                                 INT    NOT NULL,
    init_type_reply                         INT    NOT NULL,
    vl_arb_low_cap                          INT    NOT NULL,
    pkey_violations                         INT    NOT NULL,
    mkey_violations                         INT    NOT NULL,
    subn_tmo                                INT    NOT NULL,
    multicast_pkey_trap_suppression_enabled INT    NOT NULL,
    client_reregister                       INT    NOT NULL,
    guid_cap                                INT    NOT NULL,
    qkey_violations                         INT    NOT NULL,
    max_credit_hint                         INT    NOT NULL,
    overrun_errs                            INT    NOT NULL,
    local_phy_error                         INT    NOT NULL,
    resp_time_value                         INT    NOT NULL,
    link_round_trip_latency                 INT    NOT NULL,
    ooosl_mask                              TEXT   NOT NULL,
    cap_msk2                                INT,
    fec_actv                                INT,
    retrans_actv                            INT,
    UNIQUE (node_id, port_num)
);

CREATE INDEX ports_log_id_idx ON ports (log_id);

CREATE TABLE nodes_switch_info (
    node_id                 BIGINT PRIMARY KEY REFERENCES nodes (id) ON DELETE CASCADE,
    linear_fdb_cap          INT NOT NULL,
    random_fdb_cap          INT NOT NULL,
    mcast_fdb_cap           INT NOT NULL,
    linear_fdb_top          INT NOT NULL,
    def_port                INT NOT NULL,
    def_mcast_pri_port      INT NOT NULL,
    def_mcast_not_pri_port  INT NOT NULL,
    life_time_value         INT NOT NULL,
    port_state_change       INT NOT NULL,
    optimized_s_lvl_mapping INT NOT NULL,
    lids_per_port           INT NOT NULL,
    part_enf_cap            INT NOT NULL,
    inb_enf_cap             INT NOT NULL,
    outb_enf_cap            INT NOT NULL,
    filter_raw_inb_cap      INT NOT NULL,
    filter_raw_outb_cap     INT NOT NULL,
    enp0                    INT NOT NULL,
    mcast_fdb_top           INT NOT NULL
);

CREATE TABLE nodes_system_info (
    node_id       BIGINT PRIMARY KEY REFERENCES nodes (id) ON DELETE CASCADE,
    serial_number TEXT NOT NULL,
    part_number   TEXT NOT NULL,
    revision      TEXT NOT NULL,
    product_name  TEXT NOT NULL
);

CREATE TABLE nodes_sharp_info (
    node_id                   BIGINT PRIMARY KEY REFERENCES nodes (id) ON DELETE CASCADE,
    endianness                INT NOT NULL,
    enable_endianness_per_job INT NOT NULL,
    reproducibility_disable   INT NOT NULL
);
