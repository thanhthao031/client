// @flow
import React, {Component} from 'react'
import {Box, Button, HOCTimers, Icon, Text} from '../../common-adapters'
import {globalColors, globalStyles, globalMargins} from '../../styles'

import type {Props} from '.'
import type {TimerProps} from '../../common-adapters/hoc-timers'

const Splash = HOCTimers(
  class _Splash extends Component<void, Props & TimerProps, {showFeedback: boolean}> {
    timeoutId: number
    state = {
      showFeedback: false,
    }

    componentWillMount() {
      this.props.clearTimeout(this.timeoutId)
      this.timeoutId = this.props.setTimeout(() => {
        this.setState({
          showFeedback: true,
        })
      }, 4000)
    }

    render() {
      return (
        <Box style={{...stylesLoginForm, justifyContent: 'center'}}>
          <Icon type="icon-keybase-logo-80" />
          <Text style={stylesHeader} type="HeaderBig">Keybase</Text>
          {this.state.showFeedback &&
            <Box style={globalStyles.flexBoxColumn}>
              <Text style={{marginTop: globalMargins.large}} type="BodySmall">
                Keybase not starting up?
              </Text>
              <Button
                label="Let us know"
                onClick={this.props.onFeedback}
                style={{marginTop: globalMargins.small}}
                type="Primary"
              />
            </Box>}
        </Box>
      )
    }
  }
)

const Failure = (props: Props) => (
  <Box style={{...stylesLoginForm, marginTop: 0}}>
    <Box style={{...stylesBannerRed}}>
      <Text type="BodySemibold" style={stylesTextBanner}>
        Oops, we had a problem communicating with our services. This might be because you lost connectivity.
      </Text>
    </Box>
    <Box style={{...globalStyles.flexBoxColumn, flex: 1, justifyContent: 'center', alignItems: 'center'}}>
      <Icon type="icon-keybase-logo-logged-out-80" />
      <Box style={globalStyles.flexBoxRow}>
        <Button
          type="Primary"
          label="Reload"
          onClick={props.onRetry}
          style={{marginTop: globalMargins.xlarge}}
        />
        <Button
          type="Secondary"
          label="Send us feedback"
          onClick={props.onFeedback}
          style={{marginLeft: globalMargins.small, marginTop: globalMargins.xlarge}}
        />
      </Box>
    </Box>
  </Box>
)

const Intro = (props: Props) => (
  <Box
    style={{
      ...stylesLoginForm,
      marginTop: props.justRevokedSelf || props.justDeletedSelf || props.justLoginFromRevokedDevice ? 0 : 55,
    }}
  >
    {!!props.justRevokedSelf &&
      <Box style={{...stylesBannerGreen}}>
        <Text type="BodySemibold" style={stylesTextBanner}>
          <Text type="BodySemiboldItalic" style={stylesTextBanner}>{props.justRevokedSelf}</Text>
          &nbsp;was revoked successfully.
        </Text>
      </Box>}
    {!!props.justDeletedSelf &&
      <Box style={{...stylesBannerBlue}}>
        <Text type="BodySemibold" style={stylesTextBanner}>
          Your Keybase account
          {' '}
          <Text type="BodySemiboldItalic" style={stylesTextBanner}>{props.justDeletedSelf}</Text>
          &nbsp;has been deleted.
        </Text>
      </Box>}
    {!!props.justLoginFromRevokedDevice &&
      <Box style={{...stylesBannerBlue}}>
        <Text type="BodySemibold" style={stylesTextBanner}>
          Your device has been revoked, please log in again.
        </Text>
      </Box>}
    <Box
      style={{
        ...globalStyles.flexBoxColumn,
        alignItems: 'center',
        justifyContent: 'center',
        height: '100%',
        width: '100%',
      }}
    >
      <Icon type="icon-keybase-logo-80" />
      <Text style={stylesHeader} type="HeaderBig">Join Keybase</Text>
      <Text style={stylesHeaderSub} type="Body">Folders for anyone in the world.</Text>
      <Button style={stylesSignupButton} type="Primary" onClick={props.onSignup} label="Create an account" />
      <Box style={{flex: 1}} />
      <Text style={stylesLoginHeader} type="Body" onClick={props.onLogin}>Already on Keybase?</Text>
      <Button style={stylesLoginButton} type="Secondary" onClick={props.onLogin} label="Log in" />
      <Text style={stylesFeedback} type="BodySmallPrimaryLink" onClick={props.onFeedback}>
        Problems logging in?
      </Text>
    </Box>
  </Box>
)

const stylesFeedback = {
  alignSelf: 'flex-end',
  margin: globalMargins.tiny,
}

const stylesLoginForm = {
  ...globalStyles.flexBoxColumn,
  alignItems: 'center',
  flex: 1,
  justifyContent: 'flex-start',
}

const stylesHeader = {
  color: globalColors.orange,
  marginTop: globalMargins.small,
}

const stylesHeaderSub = {
  marginTop: globalMargins.tiny,
}

const stylesLoginHeader = {
  textAlign: 'center',
}

const stylesSignupButton = {
  marginTop: globalMargins.medium,
}

const stylesLoginButton = {
  margin: globalMargins.small,
}

const stylesBannerBlue = {
  ...globalStyles.flexBoxRow,
  alignItems: 'center',
  alignSelf: 'stretch',
  backgroundColor: globalColors.blue,
  justifyContent: 'center',
  marginBottom: 40,
  marginTop: 20,
  minHeight: 40,
  paddingBottom: globalMargins.tiny,
  paddingLeft: globalMargins.medium,
  paddingRight: globalMargins.medium,
  paddingTop: globalMargins.tiny,
}

const stylesBannerGreen = {
  ...stylesBannerBlue,
  backgroundColor: globalColors.green,
}

const stylesBannerRed = {
  ...stylesBannerBlue,
  backgroundColor: globalColors.red,
}

const stylesTextBanner = {
  color: globalColors.white,
  textAlign: 'center',
}

export {Intro, Failure, Splash}
