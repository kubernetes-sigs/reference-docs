<?xml version="1.0" encoding="UTF-8"?>
<xsl:stylesheet xmlns:xsl="http://www.w3.org/1999/XSL/Transform"
                version="1.0">
    <xsl:import href="/usr/share/xml/docbook/stylesheet/docbook-xsl/html/chunkfast.xsl"/>

    <xsl:param name="generate.toc">
            book      title
            part      title
    </xsl:param>
        
    <!-- UTF-8 encoding -->
    <xsl:param name="chunker.output.encoding">UTF-8</xsl:param>

    <xsl:param name="chunk.section.depth">0</xsl:param>

    <!-- add meta viewport in HEAD -->
    <xsl:template name="user.head.content">
        <meta name="viewport" content="width=device-width, user-scalable=no"/>
        <link rel="shortcut icon" type="image/png" href="images/favicon.png"/>
        <link rel="stylesheet" href="css/style.css"/>
        <link rel="stylesheet" href="css/base_fonts.css"/>
        <link rel="stylesheet" href="css/jquery-ui.min.css"/>
        <link rel="stylesheet" href="css/callouts.css"/>
        <script src="js/anchor-4.1.1.min.js"></script>
        <script src="js/jquery-3.2.1.min.js"></script>
        <script src="js/jquery-ui-1.12.1.min.js"></script>
        <script src="js/bootstrap-4.3.1.min.js"></script>
        <script src="js/sweetalert-2.1.2.min.js"></script>
        <script src="js/script.js"></script>
    </xsl:template>

    <!-- page content -->
    <xsl:template name="chunk-element-content">
        <xsl:param name="content">
            <xsl:apply-imports/>
        </xsl:param>

        <html id="docs">
            <xsl:call-template name="root.attributes"/>
            <xsl:call-template name="html.head">
            </xsl:call-template>

            <body>
                <xsl:call-template name="body.attributes"/>

                <div id="cellophane" onclick="kub.toggleMenu()"></div>
                <header>
                    <a href="/" class="logo"></a>
                    <div class="nav-buttons" data-auto-burger="primary">
                        <ul class="global-nav">
                            <li><a href="/docs/" class="active">Documentation</a></li>
                            <li><a href="/blog/">Blog</a></li>
                            <li><a href="/partners/">Partners</a></li>
                            <li><a href="/community/">Community</a></li>
                            <li><a href="/case-studies/">Case Studies</a></li>
                        </ul>
                    </div>
                </header>
                <section id="hero" class="light-text no-sub">
                    <h1>Reference</h1>
                    <h5></h5>
                    <div id="vendorStrip" class="light-text">
                        <ul>
                            <li><a href="/docs/home/">HOME</a></li>
                            <li><a href="/docs/setup/">GETTING STARTED</a></li>
                            <li><a href="/docs/concepts/">CONCEPTS</a></li>
                            <li><a href="/docs/tasks/">TASKS</a></li>
                            <li><a href="/docs/tutorials/">TUTORIALS</a></li>
                            <li><a href="/docs/reference/" class="YAH">REFERENCE</a></li>
                            <li><a href="/docs/contribute/">CONTRIBUTE</a></li>
                        </ul>
                    </div>
                </section>

                <section id="encyclopedia">
                    <div id="docsToc">                            
                            <!-- insert ToC here -->
                            <xsl:call-template name="user.header.content"/>
                            <!-- end of ToC -->
                            
                            <button class="push-menu-close-button" onclick="kub.toggleToc()"></button>
                        </div>
                        <div id="docsContent">
                            
                            <!-- Insert content here-->
                            <xsl:copy-of select="$content"/>
                            <!-- end of content -->
            
                        </div>
                    </section>
                    <footer>
                        <main class="light-text">
                            <nav>
                                <a href="/docs/home/">Home</a>
                                <a href="/blog/">Blog</a>
                                <a href="/partners/">Partners</a>
                                <a href="/community/">Community</a>
                                <a href="/case-studies/">Case Studies</a>
                            </nav>
                            <div id="miceType" class="center">(C) 2019 The Kubernetes Authors | Documentation Distributed under <a href="https://git.k8s.io/website/LICENSE" class="light-text">CC BY 4.0</a></div>
                            <div id="miceType" class="center">Copyright (C) 2019 The Linux Foundation (R) All rights reserved. The Linux Foundation has registered trademarks and uses trademarks. For a list of trademarks of The Linux Foundation, please see our <a href="https://www.linuxfoundation.org/trademark-usage" class="light-text">Trademark Usage page</a></div>
                            <div id="miceType" class="center">ICP license: 京ICP备17074266号-3</div>
                        </main>
                    </footer>
                    <button class="flyout-button" onclick="kub.toggleToc()"></button>

            </body>
        </html>
        <xsl:value-of select="$chunk.append"/>
    </xsl:template>


    <xsl:variable name="toc.listitem.type">para</xsl:variable>


    <!-- insert ToC on each chunk -*- based on make.toc template -->
    <xsl:template name="user.header.content">
        <xsl:param name="toc-context" select="/"/>
        <xsl:param name="toc.title.p" select="true()"/>
        <xsl:param name="nodes" select="/NOT-AN-ELEMENT"/>
        <xsl:variable name="root-nodes" select="/"/>

        <xsl:variable name="nodes.plus" select="$root-nodes | qandaset"/>
          
        <xsl:variable name="toc.title">
            <a class="item" data-title="Kubernetes API v1.16" href="/docs/concepts/"></a>
        </xsl:variable>
          
        <div class="pi-accordion">
                <xsl:if test="$root-nodes">
                    <xsl:copy-of select="$toc.title"/>
                        <xsl:apply-templates select="$root-nodes" mode="toc">
                            <xsl:with-param name="toc-context" select="$toc-context"/>
                        </xsl:apply-templates>
                </xsl:if>          
        </div>
    </xsl:template>


    <xsl:template name="subtoc">
        <xsl:param name="toc-context" select="."/>
        <xsl:param name="nodes" select="NOT-AN-ELEMENT"/>

        <xsl:variable name="nodes.plus" select="$nodes | qandaset"/>

        <xsl:variable name="subtoc">
                <xsl:choose>
                    <xsl:when test="$qanda.in.toc != 0">
                        <xsl:apply-templates mode="toc" select="$nodes.plus">
                            <xsl:with-param name="toc-context" select="$toc-context"/>
                        </xsl:apply-templates>
                    </xsl:when>
                    <xsl:otherwise>
                        <xsl:apply-templates mode="toc" select="$nodes">
                            <xsl:with-param name="toc-context" select="$toc-context"/>
                        </xsl:apply-templates>
                    </xsl:otherwise>
                </xsl:choose>
        </xsl:variable>

        <xsl:variable name="depth">
            <xsl:choose>
                <xsl:when test="local-name(.) = 'section'">
                    <xsl:value-of select="count(ancestor::section) + 1"/>
                </xsl:when>
                <xsl:when test="local-name(.) = 'sect1'">1</xsl:when>
                <xsl:when test="local-name(.) = 'sect2'">2</xsl:when>
                <xsl:when test="local-name(.) = 'sect3'">3</xsl:when>
                <xsl:when test="local-name(.) = 'sect4'">4</xsl:when>
                <xsl:when test="local-name(.) = 'sect5'">5</xsl:when>
                <xsl:when test="local-name(.) = 'refsect1'">1</xsl:when>
                <xsl:when test="local-name(.) = 'refsect2'">2</xsl:when>
                <xsl:when test="local-name(.) = 'refsect3'">3</xsl:when>
                <xsl:when test="local-name(.) = 'topic'">1</xsl:when>
                <xsl:when test="local-name(.) = 'simplesect'">
                    <!-- sigh... -->
                    <xsl:choose>
                        <xsl:when test="local-name(..) = 'section'">
                            <xsl:value-of select="count(ancestor::section)"/>
                        </xsl:when>
                        <xsl:when test="local-name(..) = 'sect1'">2</xsl:when>
                        <xsl:when test="local-name(..) = 'sect2'">3</xsl:when>
                        <xsl:when test="local-name(..) = 'sect3'">4</xsl:when>
                        <xsl:when test="local-name(..) = 'sect4'">5</xsl:when>
                        <xsl:when test="local-name(..) = 'sect5'">6</xsl:when>
                        <xsl:when test="local-name(..) = 'topic'">2</xsl:when>
                        <xsl:when test="local-name(..) = 'refsect1'">2</xsl:when>
                        <xsl:when test="local-name(..) = 'refsect2'">3</xsl:when>
                        <xsl:when test="local-name(..) = 'refsect3'">4</xsl:when>
                        <xsl:otherwise>1</xsl:otherwise>
                    </xsl:choose>
                </xsl:when>
                <xsl:otherwise>0</xsl:otherwise>
            </xsl:choose>
        </xsl:variable>

        <xsl:variable name="depth.from.context" select="count(ancestor::*)-count($toc-context/ancestor::*)"/>

        <xsl:variable name="subtoc.list">
            <xsl:choose>
                <xsl:when test="$toc.dd.type = ''">
                    <xsl:copy-of select="$subtoc"/>
                </xsl:when>
                <xsl:otherwise>
                    <xsl:copy-of select="$subtoc"/>
                </xsl:otherwise>
            </xsl:choose>
        </xsl:variable>

        <xsl:choose>
            <xsl:when test="$toc.listitem.type != 'li' and
                  ( (self::set or self::book or self::part) or 
                        $toc.section.depth > $depth) and 
                ( ($qanda.in.toc = 0 and count($nodes)&gt;0) or
                  ($qanda.in.toc != 0 and count($nodes.plus)&gt;0) )
                and $toc.max.depth > $depth.from.context">

                <div class="item">
                    <xsl:attribute name="data-title">
                        <xsl:apply-templates select="." mode="titleabbrev.markup"/>
                    </xsl:attribute>
                    <xsl:call-template name="toc.line">
                        <xsl:with-param name="toc-context" select="$toc-context"/>
                    </xsl:call-template>                
                    <div class="container">
                        <xsl:copy-of select="$subtoc.list"/>
                    </div>
                </div>
                        
            </xsl:when>
            <xsl:otherwise>
                <xsl:call-template name="toc.line">
                    <xsl:with-param name="toc-context" select="$toc-context"/>
                </xsl:call-template>                
            </xsl:otherwise>
        </xsl:choose>
    </xsl:template>

    <!-- Line -->
    <xsl:template name="toc.line">
        <xsl:param name="toc-context" select="."/>
        <xsl:param name="depth" select="1"/>
        <xsl:param name="depth.from.context" select="8"/>
          
          
        <a class="item">
            <xsl:attribute name="href">
                <xsl:call-template name="href.target">
                    <xsl:with-param name="context" select="$toc-context"/>
                    <xsl:with-param name="toc-context" select="$toc-context"/>
                </xsl:call-template>
            </xsl:attribute>
                        
            <!-- * if $autotoc.label.in.hyperlink is non-zero, then output the label -->
            <!-- * as part of the hyperlinked title -->
            <!--xsl:if test="not($autotoc.label.in.hyperlink = 0)">
                <xsl:variable name="label">
                <xsl:apply-templates select="." mode="label.markup"/>
                </xsl:variable>
                <xsl:copy-of select="$label"/>
                <xsl:if test="$label != ''">
                <xsl:value-of select="$autotoc.label.separator"/>
                </xsl:if>
            </xsl:if-->
            <xsl:attribute name="data-title">
                <xsl:apply-templates select="." mode="titleabbrev.markup"/>
            </xsl:attribute>
            <div class="title">
                <xsl:apply-templates select="." mode="titleabbrev.markup"/>
            </div>
        </a>
    </xsl:template>

</xsl:stylesheet>
