describe('Community Templates', () => {
  beforeEach(() => {
    cy.flush()

    cy.signin().then(({body}) => {
      const {
        org: {id},
      } = body
      cy.wrap(body.org).as('org')

      cy.fixture('routes').then(({orgs}) => {
        cy.visit(`${orgs}/${id}/settings/templates`)
      })
    })
  })

  it('The browse community template button launches github', () => {
    cy.getByTestID('browse-template-button')
      .should('have.prop', 'href')
      .and(
        'equal',
        'https://github.com/influxdata/community-templates#templates'
      )
  })

  it('The lookup template errors on invalid data', () => {
    //on empty
    cy.getByTestID('lookup-template-button').click()
    cy.getByTestID('notification-error').should('be.visible')

    //lookup template errors on bad url
    cy.getByTestID('lookup-template-input').type('www.badURL.com')
    cy.getByTestID('lookup-template-button').click()
    cy.getByTestID('notification-error').should('be.visible')

    //lookup template errors on bad file type
    cy.getByTestID('lookup-template-input').clear()
    cy.getByTestID('lookup-template-input').type('variables.html')
    cy.getByTestID('lookup-template-button').click()
    cy.getByTestID('notification-error').should('be.visible')

    //lookup template errors on github folder
    cy.getByTestID('lookup-template-input').clear()
    cy.getByTestID('lookup-template-input').type(
      'https://github.com/influxdata/community-templates/tree/master/kafka'
    )
    cy.getByTestID('lookup-template-button').click()
    cy.getByTestID('notification-error').should('be.visible')
  })

  it('Can install from CLI', () => {
    //shows file link in source
  })

  it('Simple Download', () => {
    //The lookup template accepts github raw link
    cy.getByTestID('lookup-template-input').type(
      'https://raw.githubusercontent.com/influxdata/community-templates/master/downsampling/dashboard.yml'
    )
    cy.getByTestID('lookup-template-button').click()
    cy.getByTestID('template-install-overlay').should('be.visible')

    //check that with 1 resource pluralization is correct
    cy.getByTestID('template-install-title').should('contain', 'resource')
    cy.getByTestID('template-install-title').should(
      'not.contain',
      'resources'
    )

    //check that no resources check lead to disabled install button
    cy.getByTestID('heading-Dashboards').click()
    cy.getByTestID('templates-toggle--Downsampling Status').should(
      'be.visible'
    )
    cy.getByTestID('template-install-button').should('exist')
    cy.getByTestID('templates-toggle--Downsampling Status').click()
    cy.getByTestID('template-install-button').should('not.exist')

    //and check that 0 resources pluralization is correct
    cy.getByTestID('template-install-title').should('contain', 'resources')
  })

  describe('Opening the install overlay', () => {
    beforeEach(() => {
      //lookup normal github link
      cy.getByTestID('lookup-template-input').type(
        'https://github.com/influxdata/community-templates/blob/master/docker/docker.yml'
      )
      cy.getByTestID('lookup-template-button').click()
      cy.getByTestID('template-install-overlay').should('be.visible')
    })

    it('Complicated Download', () => {
      //check that with multiple resources pluralization is correct
      cy.getByTestID('template-install-title').should('contain', 'resources')

      //no uncheck of buckets
      cy.getByTestID('template-install-title').should('contain', '22')
      cy.getByTestID('heading-Buckets').click()
      cy.getByTestID('templates-toggle--docker').should('be.visible')
      cy.getByTestID('template-install-title').should('contain', '22')
      // cy.getByTestID('templates-toggle--docker').should('be.disabled')

      //no uncheck of variables
      cy.getByTestID('template-install-title').should('contain', '22')
      cy.getByTestID('heading-Variables').click()
      cy.getByTestID('templates-toggle--bucket').should('be.visible')
      cy.getByTestID('template-install-title').should('contain', '22')
      // cy.getByTestID('templates-toggle--bucket').should('be.disabled')

      //can check and uncheck other resources
      cy.getByTestID('template-install-title').should('contain', '22')
      cy.getByTestID('heading-Checks').click()
      cy.getByTestID('templates-toggle--Container Disk Usage').should(
        'be.visible'
      )
      cy.getByTestID('templates-toggle--Container Disk Usage').click()
      cy.getByTestID('template-install-title').should('contain', '21')

      cy.getByTestID('heading-Notification Rules').click()
      cy.getByTestID('templates-toggle--Crit Notifier').should('be.visible')
      cy.getByTestID('templates-toggle--Crit Notifier').click()
      cy.getByTestID('template-install-title').should('contain', '20')
    })

    it('Can install template', () => {
      cy.getByTestID('template-install-button').click()
      cy.getByTestID('notification-success').should('be.visible')
      cy.getByTestID('installed-template-docker').should('be.visible')
    })
  })

  describe('Install Completed', () => {
    beforeEach(() => {
      cy.getByTestID('lookup-template-input').type(
        'https://github.com/influxdata/community-templates/blob/master/docker/docker.yml'
      )
      cy.getByTestID('lookup-template-button').click()
      cy.getByTestID('template-install-overlay').should('be.visible')
      cy.getByTestID('template-install-button').should('exist')
      cy.getByTestID('template-install-button').click()
      cy.getByTestID('notification-success').should('be.visible')
      cy.getByTestID('installed-template-docker').should('be.visible')
    })

    it('Install Identical template', () => {
      cy.getByTestID('lookup-template-input').clear()
      cy.getByTestID('lookup-template-input').type(
        'https://github.com/influxdata/community-templates/blob/master/docker/docker.yml'
      )
      cy.getByTestID('lookup-template-button').click()
      cy.getByTestID('template-install-overlay').should('be.visible')
      cy.getByTestID('template-install-button').should('exist')
      cy.getByTestID('template-install-button').click()
      cy.getByTestID('notification-success').should('be.visible')
      cy.getByTestID('installed-template-list').should('have', '2')
    })

    it('Can click on template resources', () => {
      //button
      // cy.getByTestID('template-resource-link-Buckets')
    })

    it('Click on source takes you to github', () => {
      cy.getByTestID('template-source-link').should(
        'contain',
        'https://github.com/influxdata/community-templates/blob/master/docker/docker.yml'
      )
      //TODO: add the link from CLI
    })

    it('Can delete template', () => {
      cy.getByTestID('template-delete-button-docker--button').click()
      cy.getByTestID('template-delete-button-docker--confirm-button').click()
      cy.getByTestID('installed-template-docker').should('not.be.visible')
    })
  })
})
