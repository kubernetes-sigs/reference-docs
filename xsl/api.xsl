<?xml version="1.0" encoding="UTF-8"?>
<xsl:stylesheet xmlns:xsl="http://www.w3.org/1999/XSL/Transform"
                xmlns:fo="http://www.w3.org/1999/XSL/Format"
                version="1.0">
  <xsl:import href="/usr/share/xml/docbook/stylesheet/docbook-xsl/fo/docbook.xsl"/>

  <xsl:param name="paper.type">A4</xsl:param>
  <xsl:param name="double.sided">1</xsl:param>
  <xsl:param name="variablelist.as.blocks">1</xsl:param>
  <xsl:param name="section.autolabel">1</xsl:param>
  <xsl:param name="generate.toc">
    book      toc,title
    part      title
  </xsl:param>
  <xsl:attribute-set name="monospace.verbatim.properties">
    <xsl:attribute name="wrap-option">wrap</xsl:attribute>
  </xsl:attribute-set>
  <xsl:param name="insert.link.page.number">yes</xsl:param>
</xsl:stylesheet>
